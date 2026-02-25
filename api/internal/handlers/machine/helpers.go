package machine

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/requestx"
	"booksonhooks.ca/internal/validator"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *MachineHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to write json: %v", err)
	}
}

func mapLoadBookToBookMachine(machineID int64, book domain.BookMachine) domain.BookMachine {
	return domain.BookMachine{
		MachineID: machineID,
		BookID:    book.BookID,
		Row:       book.Row,
		Col:       book.Col,
	}
}

func mapMachineUpsertToMachine(id int64, req domain.MachineUpsert) *domain.Machine {
	return &domain.Machine{
		ID:       id,
		Location: req.Location,
		Rows:     req.Rows,
		Cols:     req.Cols,
	}
}

func mapMachineRequestToMachine(req domain.MachineRequest) *domain.Machine {
	return &domain.Machine{
		Location: req.Location,
		Rows:     req.Rows,
		Cols:     req.Cols,
	}
}

func validateMachineUpsertFields(location string, rows, cols int) map[string]string {
	fieldErrors := make(map[string]string)

	if !validator.NotBlank(location) {
		fieldErrors["location"] = "This field cannot be blank"
	} else if !validator.MaxChars(location, 100) {
		fieldErrors["location"] = "This field can have a maximum of 100 characters"
	}

	if !validator.PositiveInt(rows) {
		fieldErrors["rows"] = "Must be a positive integer"
	}

	if !validator.PositiveInt(cols) {
		fieldErrors["cols"] = "Must be a positive integer"
	}

	return fieldErrors
}

type slot struct {
	row int
	col int
}

func validateAndMapLoadBooks(machineID int64, machine *domain.Machine, payloadBooks []domain.BookMachine) ([]domain.BookMachine, *httpErrors.HTTPError) {
	seenBookIDs := make(map[int64]struct{}, len(payloadBooks))
	seenSlots := make(map[slot]struct{}, len(payloadBooks))

	books := make([]domain.BookMachine, len(payloadBooks))
	for i, b := range payloadBooks {
		if b.BookID <= 0 || b.Row < 0 || b.Col < 0 {
			return nil, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("invalid book or slot values")}
		}
		if b.Row >= machine.Rows || b.Col >= machine.Cols {
			return nil, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("slot out of bounds")}
		}
		if _, exists := seenBookIDs[b.BookID]; exists {
			return nil, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("duplicate book_id in payload")}
		}
		seenBookIDs[b.BookID] = struct{}{}

		key := slot{row: b.Row, col: b.Col}
		if _, exists := seenSlots[key]; exists {
			return nil, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("duplicate slot in payload")}
		}
		seenSlots[key] = struct{}{}

		books[i] = mapLoadBookToBookMachine(machineID, b)
	}

	return books, nil
}

func (h *MachineHandler) recordMachineMetric(r *http.Request, machineID int64) {
	var (
		isQr      bool
		source    string
		sessionID string
		hasErrors bool
	)

	parsedQr, err := requestx.ParseOptionalQueryBool(r, "is_qr")
	if err != nil {
		h.log.Warn("skipping machine metric: invalid is_qr machine_id=%d raw=%q err=%v",
			machineID, requestx.QueryString(r, "is_qr"), err)
		hasErrors = true
	} else {
		isQr = parsedQr
	}

	parsedSource, err := requestx.ParseOptionalQuerySource(r, "source")
	if err != nil {
		h.log.Warn("skipping machine metric: invalid source machine_id=%d raw=%q err=%v",
			machineID, requestx.QueryString(r, "source"), err)
		hasErrors = true
	} else {
		source = parsedSource
	}

	parsedSessionID, err := requestx.ParseOptionalQuerySessionID(r, "session_id")
	if err != nil {
		h.log.Warn("skipping machine metric: invalid session_id machine_id=%d raw=%q err=%v",
			machineID, requestx.QueryString(r, "session_id"), err)
		hasErrors = true
	} else {
		sessionID = parsedSessionID
	}

	if hasErrors {
		return
	}

	admin := requestx.IsAdminSource(source)

	if _, err := h.repo.InsertMachineMetric(r.Context(), machineID, isQr, source, admin, sessionID); err != nil {
		h.log.Error("failed to insert machine metric: machine_id=%d is_qr=%t source=%q admin=%t session_id=%q err=%v",
			machineID, isQr, source, admin, sessionID, err)
	}
}

package machine

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/httpErrors"
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

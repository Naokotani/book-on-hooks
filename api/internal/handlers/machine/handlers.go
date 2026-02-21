package machine

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/jackc/pgx/v5"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/validator"
)

type MachineRepo interface {
	GetMachines(ctx context.Context) ([]domain.Machine, error)
	GetMachineById(ctx context.Context, id int64) (*domain.Machine, error)
	GetMachineWithBooks(ctx context.Context, id int64) (*domain.MachineWithBooks, error)
	InsertMachine(ctx context.Context, machine *domain.Machine) (int64, error)
	InsertMachineMetric(ctx context.Context, machineID int64, qr bool, source string, admin bool) (int64, error)
	LoadMachine(ctx context.Context, machineID int64, books []domain.BookMachine) error
	UpdateMachine(ctx context.Context, machine *domain.Machine) error
	DeleteMachine(ctx context.Context, id int64) error
	ClearMachineBooks(ctx context.Context, machineID int64) error
}

type MachineHandler struct {
	log  *logger.Logger
	repo MachineRepo
}

func New(log *logger.Logger,
	repo MachineRepo) *MachineHandler {

	return &MachineHandler{
		log:  log,
		repo: repo,
	}
}

func (h *MachineHandler) CreateMachine(w http.ResponseWriter, r *http.Request) {
	var payload domain.MachineRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Error("failed to decode machine json: %v", err)
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	payload.Location = strings.TrimSpace(payload.Location)
	if !validator.NotBlank(payload.Location) || !validator.MaxChars(payload.Location, 100) {
		httpErrors.ClientError(w, http.StatusUnprocessableEntity)
		return
	}
	if !validator.PositiveInt(payload.Rows) || !validator.PositiveInt(payload.Cols) {
		httpErrors.ClientError(w, http.StatusUnprocessableEntity)
		return
	}

	machine := mapMachineRequestToMachine(payload)

	id, err := h.repo.InsertMachine(r.Context(), machine)
	if err != nil {
		h.log.Error("failed to insert machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	machine.ID = id
	h.writeJSON(w, http.StatusCreated, machine)
}

func (h *MachineHandler) LoadMachine(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	machine, err := h.repo.GetMachineById(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to get machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	var payload domain.MachineLoadRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Error("failed to decode load machine json: %v", err)
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	books, httpErr := validateAndMapLoadBooks(int64(id), machine, payload.Books)
	if httpErr != nil {
		httpErrors.ClientError(w, httpErr.Status)
		return
	}

	if err := h.repo.LoadMachine(r.Context(), int64(id), books); err != nil {
		h.log.Error("failed to load machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, domain.MachineLoadResponse{
		MachineID: int64(id),
		Count:     len(books),
	})
}

func (h *MachineHandler) GetMachine(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	machine, err := h.repo.GetMachineById(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to get machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, machine)
}

func (h *MachineHandler) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	var payload domain.MachineUpsert
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	payload.Location = strings.TrimSpace(payload.Location)
	if !validator.NotBlank(payload.Location) || !validator.MaxChars(payload.Location, 100) {
		httpErrors.ClientError(w, http.StatusUnprocessableEntity)
		return
	}
	if !validator.PositiveInt(payload.Rows) || !validator.PositiveInt(payload.Cols) {
		httpErrors.ClientError(w, http.StatusUnprocessableEntity)
		return
	}

	machine := mapMachineUpsertToMachine(int64(id), payload)

	if err := h.repo.UpdateMachine(r.Context(), machine); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to update machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, machine)
}

func (h *MachineHandler) DeleteMachine(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	if _, err := h.repo.GetMachineById(r.Context(), int64(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to get machine before delete: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	if err := h.repo.DeleteMachine(r.Context(), int64(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to delete machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MachineHandler) ClearMachineBooks(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	if err := h.repo.ClearMachineBooks(r.Context(), int64(id)); err != nil {
		h.log.Error("failed to clear machine books: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MachineHandler) GetMachineWithBooks(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	result, err := h.repo.GetMachineWithBooks(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to get machine with books: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.recordMachineMetric(r, int64(id))

	h.writeJSON(w, http.StatusOK, result)
}

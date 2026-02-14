package machine

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/template"
	"booksonhooks.ca/internal/validator"
)

type MachineRepo interface {
	GetMachines(ctx context.Context) ([]domain.Machine, error)
	GetMachineById(ctx context.Context, id int64) (*domain.Machine, error)
	InsertMachine(ctx context.Context, machine *domain.Machine) (int64, error)
}

type MachineHandler struct {
	temp *template.Templates
	form *forms.Form
	log  *logger.Logger
	repo MachineRepo
}

type MachineRow struct {
	books []domain.Book
}

func New(temp *template.Templates,
	form *forms.Form,
	log *logger.Logger,
	repo MachineRepo) *MachineHandler {

	return &MachineHandler{
		temp: temp,
		form: form,
		log:  log,
		repo: repo,
	}
}

func (h *MachineHandler) MachinesView(w http.ResponseWriter, r *http.Request) {
	machines, err := h.repo.GetMachines(r.Context())

	if err != nil {
		h.log.Error("Failed to read machines. %s", err)
		httpErrors.NotFound(w)
		return
	}

	h.writeJSON(w, http.StatusOK, machines)
}

func (h *MachineHandler) MachineCreateView(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	data.Form = forms.MachineCreateForm{}
	h.temp.Render(w, http.StatusOK, "newMachine.gotmpl", data)
}

func (h *MachineHandler) MachineCreate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Location string `json:"location"`
		Rows     int    `json:"rows"`
		Cols     int    `json:"cols"`
	}

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

	machine := &domain.Machine{
		Location: payload.Location,
		Rows:     payload.Rows,
		Columns:  payload.Cols,
	}

	id, err := h.repo.InsertMachine(r.Context(), machine)
	if err != nil {
		h.log.Error("failed to insert machine: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusCreated, struct {
		ID int64 `json:"id"`
	}{ID: id})
}

func (h *MachineHandler) MachineLoadView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	machine, err := h.repo.GetMachineById(r.Context(), int64(id))
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	response := struct {
		ID       int64  `json:"id"`
		Location string `json:"location"`
		Rows     int    `json:"rows"`
		Cols     int    `json:"cols"`
	}{
		ID:       machine.ID,
		Location: machine.Location,
		Rows:     machine.Rows,
		Cols:     machine.Columns,
	}

	h.writeJSON(w, http.StatusOK, response)
}

package machine

import (
	"net/http"

	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/template"
)

type MachineHandler struct {
	temp *template.Templates
	form *forms.Form
	log  *logger.Logger
	repo *repository.Database
}

type MachineRow struct {
	books []repository.Book
}

func New(temp *template.Templates,
	form *forms.Form,
	log *logger.Logger,
	repo *repository.Database) *MachineHandler {

	return &MachineHandler{
		temp: temp,
		form: form,
		log:  log,
		repo: repo,
	}
}

func (h *MachineHandler) MachinesView(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)

	machines, err := h.repo.GetMachines(r.Context())

	//TODO make error handler. Need to type check to see if its just empty.
	if err != nil {
		h.log.Error("Failed to get machines")
	}

	data.Machines = machines

	h.temp.Render(w, http.StatusOK, "machines.gotmpl", data)
}

func (h *MachineHandler) MachineCreateView(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	data.Form = forms.MachineCreateForm{}
	h.temp.Render(w, http.StatusOK, "newMachine.gotmpl", data)
}

func (h *MachineHandler) MachineCreate(w http.ResponseWriter, r *http.Request) {
	_, form, httpErr := forms.MachineFormService(h.form, r)

	if httpErr != nil {
		data := h.temp.NewTemplateData(r)
		data.Form = form
		h.temp.Render(w, httpErr.Status, "newMachine.gotmpl", data)
		return
	}

	http.Redirect(w, r, "/api/machines", http.StatusSeeOther)
}

func (h *MachineHandler) MachineLoadView(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)

	h.temp.Render(w, http.StatusOK, "loadMachine.gotmpl", data)
}

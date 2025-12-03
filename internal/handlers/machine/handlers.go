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

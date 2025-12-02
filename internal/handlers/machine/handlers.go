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

func (b *MachineHandler) MachinesView(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "machines.gotmpl", data)
}

func (b *MachineHandler) MachineCreateView(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "machines.gotmpl", data)
}

func (b *MachineHandler) MachineCreate(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "machines.gotmpl", data)
}

package home

import (
	"net/http"

	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/template"
)

type HomeHandler struct {
	temp *template.Templates
	form *forms.Form
	log  *logger.Logger
	repo *repository.Database
}

func New(temp *template.Templates,
	form *forms.Form,
	log *logger.Logger,
	repo *repository.Database) *HomeHandler {

	return &HomeHandler{
		temp: temp,
		form: form,
		log:  log,
		repo: repo,
	}
}

func (b *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "home.gotmpl", data)
}

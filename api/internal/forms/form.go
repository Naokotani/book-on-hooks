package forms

import (
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"github.com/go-playground/form/v4"
)

type Form struct {
	formDecoder *form.Decoder
	repo        *repository.Database
	log         *logger.Logger
}

func New(repo *repository.Database, log *logger.Logger) *Form {
	formDecoder := form.NewDecoder()

	return &Form{
		log:         log,
		repo:        repo,
		formDecoder: formDecoder,
	}
}

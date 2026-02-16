package app

import (
	"net/http"

	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/handlers/book"
	"booksonhooks.ca/internal/handlers/machine"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
)

type Application struct {
	log *logger.Logger
	db  *repository.Database

	machineHandlers *machine.MachineHandler
	bookHandlers    *book.BookHandler
}

func CreateApp(addr *string, logger *logger.Logger, db *repository.Database) http.Server {
	app := &Application{
		log: logger,
		db:  db,
	}

	form := forms.New(db, logger)
	app.bookHandlers = book.New(form, logger, db)
	app.machineHandlers = machine.New(logger, db)

	return http.Server{
		Addr:     *addr,
		ErrorLog: logger.ErrorLog,
		Handler:  app.Routes(),
	}
}

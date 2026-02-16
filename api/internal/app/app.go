package app

import (
	"net/http"

	"booksonhooks.ca/internal/template"

	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/handlers/book"
	"booksonhooks.ca/internal/handlers/home"
	"booksonhooks.ca/internal/handlers/machine"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
)

type Application struct {
	log *logger.Logger
	db  *repository.Database

	machineHandlers *machine.MachineHandler
	homeHandlers    *home.HomeHandler
	bookHandlers    *book.BookHandler
}

func CreateApp(addr *string, logger *logger.Logger, db *repository.Database) http.Server {
	app := &Application{
		log: logger,
		db:  db,
	}

	form := forms.New(db, logger)
	template, err := template.New(app.ServerError)

	if err != nil {
		logger.ErrorLog.Printf("Failed to create templates\n%s", err)
	}

	app.homeHandlers = home.New(template, form, logger, db)
	app.bookHandlers = book.New(template, form, logger, db)
	app.machineHandlers = machine.New(template, form, logger, db)

	return http.Server{
		Addr:     *addr,
		ErrorLog: logger.ErrorLog,
		Handler:  app.Routes(),
	}
}

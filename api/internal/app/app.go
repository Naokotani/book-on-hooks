package app

import (
	"net/http"

	"booksonhooks.ca/internal/handlers/adminauth"
	"booksonhooks.ca/internal/handlers/book"
	"booksonhooks.ca/internal/handlers/health"
	"booksonhooks.ca/internal/handlers/machine"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"github.com/alexedwards/scs/v2"
)

type Application struct {
	log *logger.Logger
	db  *repository.Database

	sessions *scs.SessionManager

	machineHandlers *machine.MachineHandler
	bookHandlers    *book.BookHandler
	healthHandlers  *health.HealthHandler
	authHandlers    *adminauth.Handler
}

func CreateApp(addr *string, logger *logger.Logger, db *repository.Database) http.Server {
	app := &Application{
		log:      logger,
		db:       db,
		sessions: newSessionManager(db.Db),
	}

	app.bookHandlers = book.New(logger, db)
	app.machineHandlers = machine.New(logger, db)
	app.healthHandlers = health.New(logger, db)
	app.authHandlers = adminauth.New(logger, app.sessions)

	return http.Server{
		Addr:     *addr,
		ErrorLog: logger.ErrorLog,
		Handler:  app.Routes(),
	}
}

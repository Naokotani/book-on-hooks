package main

import (
	//"github.com/Naokotani/logger"
	"booksonhooks.ca/internal/logger"
)

type application struct {
	log *logger.Logger
}

func main() {
	logger := logger.NewLogger("info")
	app := application{
		log: &logger,
	}

	const port = ":4000"
	app.log.Info("Starting application on port %s\n", port)
}

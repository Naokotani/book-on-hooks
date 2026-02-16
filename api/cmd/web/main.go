package main

import (
	"flag"

	"booksonhooks.ca/internal/app"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
)

func main() {
	logger := logger.NewLogger("info")

	const port = ":4000"
	logger.Info("Starting application on port %s", port)

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	db, err := repository.GetDatabaseConnection()
	if err != nil {
		logger.ErrorLog.Fatalf("failed to create database connection\n%s", err)
	}
	defer db.Db.Close()

	srv := app.CreateApp(addr, &logger, db)
	err = srv.ListenAndServe()
	logger.ErrorLog.Fatal(err)
}

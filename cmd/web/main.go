package main

import (
	"flag"

	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"html/template"
	"net/http"
)

type application struct {
	log           *logger.Logger
	templateCache map[string]*template.Template
	db            repository.Database
}

func main() {
	logger := logger.NewLogger("info")

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

	db, err := repository.GetDatabaseConnection()

	app := application{
		log:           &logger,
		templateCache: templateCache,
		db:            db,
	}

	err = app.db.CreateTables()

	if err != nil {
		app.log.ErrorLog.Fatalf("Failed to create database tables.\n%s\n", err)
	}

	book := &repository.Book{Title: "Alice in wonderland", Price: "12.99", Cover: "foo", Author: "lewis carol"}
	err = app.db.InsertBook(book)

	const port = ":4000"
	app.log.Info("Starting application on port %s\n", port)

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: logger.ErrorLog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	logger.ErrorLog.Fatal(err)
}

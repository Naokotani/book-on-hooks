package main

import (
	"html/template"
	"path/filepath"
	"time"

	"booksonhooks.ca/internal/repository"
)

type templateData struct {
	CurrentYear int
	Row         int
	Col         int
	Form        any
	Book        *repository.Book
}

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.gotmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gotmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.gotmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

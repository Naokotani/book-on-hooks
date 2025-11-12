package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	book, err := app.db.GetBook(1)

	if err != nil {
		app.log.Error("Failed to retrieve book.\n%s\n", err)
	}

	app.log.Info("Book title: %s", book.Title)

	app.render(w, http.StatusOK, "home.gotmpl", data)
}

func (app *application) covers(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	row, err := strconv.Atoi(params.ByName("row"))
	if err != nil || row < 0 {
		app.notFound(w)
		return
	}

	col, err := strconv.Atoi(params.ByName("col"))
	if err != nil || col < 0 {
		app.notFound(w)
		return
	}

	app.log.Info("retrieve cover at row: %d, col: %d", row, col)

	data := app.newTemplateData(r)
	data.Col = col
	data.Row = row

	app.render(w, http.StatusOK, "covers.gotmpl", data)
}

func (app *application) book(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)

}
func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)
}

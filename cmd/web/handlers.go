package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "home.gotmpl", data)
}

func (app *application) cover(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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

	book, err := app.db.GetBookByRowAndCol(row, col)
	if err != nil || col < 0 {
		app.notFound(w)
		return
	}

	data.Book = book

	app.log.Info("retrieve cover at row: %d, col: %d", row, col)

	data.Col = 1
	data.Row = 1

	app.render(w, http.StatusOK, "covers.gotmpl", data)
}

func (app *application) book(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)
}

func (app *application) getBooks(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "books.gotmpl", data)
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.log.Error("Failed to parse form: %v", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, form, httpErr := parseBookForm(app, r)

	if httpErr.Err != nil {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, httpErr.Status, "newBook.gotmpl", data)
		return
	}

	_, err = app.db.GetBookByID(id)

	if err != nil {
		app.log.Error("Failed to retrieve newly created book: %v", err)
		app.clientError(w, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("covers/%d/%d/%d", 1, 1, 1), http.StatusSeeOther)
}

func (app *application) bookCreateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = bookCreateForm{}
	app.render(w, http.StatusOK, "newBook.gotmpl", data)
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "book.gotmpl", data)
}

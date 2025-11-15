package main

import (
	"net/http"
	"strconv"

	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type bookCreateForm struct {
	Title               string `form:"title"`
	Author              string `form:"author"`
	Price               string `form:"price"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/covers/1/1", http.StatusSeeOther)
}

func (app *application) covers(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())
	book, err := app.db.GetBookByID(1)

	if err != nil {
		app.log.Error("Failed to retrieve book.\n%s\n", err)
	}

	app.log.Info("Book before row/col, %s", book.Title)

	app.db.UpdateBookPosition(28, 4, 4)

	if err != nil {
		app.log.Error("Failed to retrieve book.\n%s\n", err)
	}

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

	book, err = app.db.GetBookByRowAndCol(row, col)
	if err != nil || col < 0 {
		app.notFound(w)
		return
	}

	data.Book = book

	app.log.Info("retrieve cover at row: %d, col: %d", row, col)

	data.Col = book.Col
	data.Row = book.Row

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
		app.log.Error("Failed to parse form: %s\n", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form bookCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.log.Error("Failed to decode form: %s\n", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.log.Info("Title: %s", form.Title)
	app.log.Info("Author: %s", form.Author)

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field can have a maximum of 100 characters")
	form.CheckField(validator.NotBlank(form.Author), "author", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Author, 100), "author", "This field can have a maximum of 100 characters")
	form.CheckField(validator.ValidDollarValue(form.Price), "price", "Must be a valid dolloar value. IE 10.99 or $10.99")

	if !form.Valid() {
		app.log.Warn("Cannot process form.")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "newBook.gotmpl", data)
		return
	}

	book := repository.Book{Title: form.Title, Author: form.Author, Price: form.Price}

	file, header, err := r.FormFile("cover")
	if err != nil {
		app.log.Error("Failed to upload file: %s\n", err)
	}
	defer file.Close()

	_, err = app.db.InsertBook(&book, file, header)

	if err != nil {
		form.CheckField(false, "cover", "Failed to upload file")
		data := app.newTemplateData(r)
		data.Form = form
		app.log.Error("Failed to insert book: %s\n", err)
		app.render(w, http.StatusUnprocessableEntity, "newBook.gotmpl", data)
		return
	}

	http.Redirect(w, r, "/covers/1/1", http.StatusSeeOther)
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

package book

import (
	"fmt"
	"net/http"
	"strconv"

	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/template"
	"github.com/julienschmidt/httprouter"
)

type BookHandler struct {
	temp *template.Templates
	form *forms.Form
	log  *logger.Logger
	repo *repository.Database
}

func New(temp *template.Templates,
	form *forms.Form,
	log *logger.Logger,
	repo *repository.Database) *BookHandler {

	return &BookHandler{
		temp: temp,
		form: form,
		log:  log,
		repo: repo,
	}
}

func (b *BookHandler) Cover(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	params := httprouter.ParamsFromContext(r.Context())

	row, err := strconv.Atoi(params.ByName("row"))
	if err != nil || row < 0 {
		httpErrors.NotFound(w)
		return
	}

	col, err := strconv.Atoi(params.ByName("col"))
	if err != nil || col < 0 {
		httpErrors.NotFound(w)
		return
	}

	book, err := b.repo.GetBookByRowAndCol(row, col)
	if err != nil || col < 0 {
		httpErrors.NotFound(w)
		return
	}

	data.Book = book

	b.log.Info("retrieve cover at row: %d, col: %d", row, col)

	data.Col = 1
	data.Row = 1

	b.temp.Render(w, http.StatusOK, "covers.gotmpl", data)
}

func (b *BookHandler) Book(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

func (b *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "books.gotmpl", data)
}

func (b *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

func (b *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		b.log.Error("Failed to parse form: %v", err)
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	id, form, httpErr := forms.ParseBookForm(b.form, r)

	if httpErr.Err != nil {
		data := b.temp.NewTemplateData(r)
		data.Form = form
		b.temp.Render(w, httpErr.Status, "newBook.gotmpl", data)
		return
	}

	_, err = b.repo.GetBookByID(id)

	if err != nil {
		b.log.Error("Failed to retrieve newly created book: %v", err)
		httpErrors.ClientError(w, http.StatusNotFound)
		httpErrors.ClientError(w, http.StatusNotFound)
	}

	http.Redirect(w, r, fmt.Sprintf("covers/%d/%d/%d", 1, 1, 1), http.StatusSeeOther)
}

func (b *BookHandler) BookCreateView(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	data.Form = forms.BookCreateForm{}
	b.temp.Render(w, http.StatusOK, "newBook.gotmpl", data)
}

func (b *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	data := b.temp.NewTemplateData(r)
	b.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

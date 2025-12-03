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

func (h *BookHandler) Cover(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
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

	book, err := h.repo.GetBookByRowAndCol(row, col)
	if err != nil || col < 0 {
		httpErrors.NotFound(w)
		return
	}

	data.Book = book

	h.log.Info("retrieve cover at row: %d, col: %d", row, col)

	data.Col = 1
	data.Row = 1

	h.temp.Render(w, http.StatusOK, "covers.gotmpl", data)
}

func (h *BookHandler) Book(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	h.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	h.temp.Render(w, http.StatusOK, "books.gotmpl", data)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	h.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.log.Error("Failed to parse form: %v", err)
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	id, form, httpErr := forms.ParseBookForm(h.form, r)

	if httpErr.Err != nil {
		data := h.temp.NewTemplateData(r)
		data.Form = form
		h.temp.Render(w, httpErr.Status, "newBook.gotmpl", data)
		return
	}

	_, err = h.repo.GetBookByID(id)

	if err != nil {
		h.log.Error("Failed to retrieve newly created book: %v", err)
		httpErrors.ClientError(w, http.StatusNotFound)
	}

	http.Redirect(w, r, fmt.Sprintf("covers/%d/%d/%d", 1, 1, 1), http.StatusSeeOther)
}

func (h *BookHandler) BookCreateView(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	data.Form = forms.BookCreateForm{}
	h.temp.Render(w, http.StatusOK, "newBook.gotmpl", data)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	h.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

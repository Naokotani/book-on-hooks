package book

import (
	"encoding/json"
	"net/http"
	"strconv"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/template"
	"context"
	"github.com/julienschmidt/httprouter"
)

type BookRepo interface {
	GetBookByRowAndCol(ctx context.Context, row, col int) (*domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	GetBookByID(ctx context.Context, id int64) (*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	DeleteBook(ctx context.Context, id int64) error
}

type BookHandler struct {
	temp *template.Templates
	form *forms.Form
	log  *logger.Logger
	repo BookRepo
}

func New(temp *template.Templates,
	form *forms.Form,
	log *logger.Logger,
	repo BookRepo) *BookHandler {

	return &BookHandler{
		temp: temp,
		form: form,
		log:  log,
		repo: repo,
	}
}

func (h *BookHandler) Cover(w http.ResponseWriter, r *http.Request) {
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

	book, err := h.repo.GetBookByRowAndCol(r.Context(), row, col)
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	h.log.Info("retrieve cover at row: %d, col: %d", row, col)

	response := struct {
		Row  int `json:"row"`
		Col  int `json:"col"`
		Book any `json:"book"`
	}{
		Row:  row,
		Col:  col,
		Book: book,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *BookHandler) Book(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	h.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.GetBooks(r.Context())

	//TODO make error handler. Need to type check to see if its just empty.
	if err != nil {
		httpErrors.NotFound(w)
		h.log.Error("Failed to read books. %s", err)
		return
	}

	for _, b := range books {
		h.log.Info("Book: %v", b)
	}

	h.writeJSON(w, http.StatusOK, books)
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

	id, form, httpErr := forms.BookFormService(h.form, r)

	if httpErr != nil {
		data := h.temp.NewTemplateData(r)
		data.Form = form
		h.temp.Render(w, httpErr.Status, "newBook.gotmpl", data)
		return
	}

	_, err = h.repo.GetBookByID(r.Context(), id)

	if err != nil {
		h.log.Error("Failed to retrieve newly created book: %v", err)
		httpErrors.ClientError(w, http.StatusNotFound)
	}

	http.Redirect(w, r, "/api/books", http.StatusSeeOther)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	data := h.temp.NewTemplateData(r)
	h.temp.Render(w, http.StatusOK, "book.gotmpl", data)
}

func (h *BookHandler) GetBookByIDJSON(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	book, err := h.repo.GetBookByID(r.Context(), int64(id))
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	var payload struct {
		Title   string `json:"title"`
		Author  string `json:"author"`
		Summary string `json:"summary"`
		Price   string `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	book := &domain.Book{
		ID:      int64(id),
		Title:   payload.Title,
		Author:  payload.Author,
		Summary: payload.Summary,
		Price:   payload.Price,
	}

	if err := h.repo.UpdateBook(r.Context(), book); err != nil {
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id <= 0 {
		httpErrors.NotFound(w)
		return
	}

	if err := h.repo.DeleteBook(r.Context(), int64(id)); err != nil {
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

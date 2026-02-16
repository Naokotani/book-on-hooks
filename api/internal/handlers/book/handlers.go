package book

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/forms"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/logger"
	"github.com/julienschmidt/httprouter"
)

type BookRepo interface {
	GetBooks(ctx context.Context) ([]domain.Book, error)
	GetBookByID(ctx context.Context, id int64) (*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	DeleteBook(ctx context.Context, id int64) error
}

type BookHandler struct {
	form *forms.Form
	log  *logger.Logger
	repo BookRepo
}

func New(form *forms.Form,
	log *logger.Logger,
	repo BookRepo) *BookHandler {

	return &BookHandler{
		form: form,
		log:  log,
		repo: repo,
	}
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

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.log.Error("Failed to parse form: %v", err)
		h.writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "failed to parse multipart form",
		})
		return
	}

	id, form, httpErr := forms.BookFormService(h.form, r)

	if httpErr != nil {
		response := map[string]any{
			"error": httpErr.Error(),
		}
		if form != nil && len(form.FieldErrors) > 0 {
			response["field_errors"] = form.FieldErrors
		}
		h.writeJSON(w, httpErr.Status, response)
		return
	}

	book, err := h.repo.GetBookByID(r.Context(), id)

	if err != nil {
		h.log.Error("Failed to retrieve newly created book: %v", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "book was created but could not be retrieved",
			"id":    id,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, book)
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

package book

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/requestx"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

const defaultImageDir = "/data/images"

type BookRepo interface {
	GetBooks(ctx context.Context) ([]domain.Book, error)
	GetBookByID(ctx context.Context, id int64) (*domain.Book, error)
	GetBookLocations(ctx context.Context, bookID int64) (*domain.BookLocation, error)
	InsertBook(ctx context.Context, book *domain.Book, file multipart.File, header *multipart.FileHeader) (int64, error)
	UpdateBook(ctx context.Context, book *domain.Book) error
	UpdateBookImage(ctx context.Context, id int64, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteBook(ctx context.Context, id int64) (string, error)
	InsertBookMetric(ctx context.Context, bookID, machineID int64, qr bool, source, sessionID string) (int64, error)
}

type BookHandler struct {
	log  *logger.Logger
	repo BookRepo
}

func New(log *logger.Logger,
	repo BookRepo) *BookHandler {

	return &BookHandler{
		log:  log,
		repo: repo,
	}
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.GetBooks(r.Context())

	if err != nil {
		h.log.Error("failed to read books: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := requestx.ParsePathID(r, "id")
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	book, err := h.repo.GetBookByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to get book: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) GetBookSummary(w http.ResponseWriter, r *http.Request) {
	id, err := requestx.ParsePathID(r, "id")
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	bookLocation, err := h.repo.GetBookLocations(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to get book summary: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.recordBookMetric(r, id)

	h.writeJSON(w, http.StatusOK, bookLocation)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.log.Error("failed to parse form: %v", err)
		h.writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "failed to parse multipart form",
		})
		return
	}

	createReq, fieldErrors := validateCreateBookMetadata(r)
	if len(fieldErrors) > 0 {
		h.log.Warn("invalid create-book payload: field_errors=%v", fieldErrors)
		httpErrors.ValidationError(w, fieldErrors)
		return
	}

	file, header, fieldErrors := validateBookImage(r)
	if len(fieldErrors) > 0 {
		h.log.Warn("invalid create-book image: field_errors=%v", fieldErrors)
		httpErrors.ValidationError(w, fieldErrors)
		return
	}
	defer file.Close()

	id, err := h.repo.InsertBook(r.Context(), mapBookCreateRequestToBook(createReq), file, header)
	if err != nil {
		h.log.Error("failed to insert book: %v", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to insert book",
		})
		return
	}

	book, err := h.repo.GetBookByID(r.Context(), id)

	if err != nil {
		h.log.Error("failed to retrieve newly created book: %v", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "book was created but could not be retrieved",
			"id":    id,
		})
		return
	}

	h.writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := requestx.ParsePathID(r, "id")
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	var payload domain.BookUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	metadata := mapBookUpdateRequestToMetadata(payload)
	if fieldErrors := validateBookMetadata(metadata); len(fieldErrors) > 0 {
		h.log.Warn("invalid update-book payload: field_errors=%v", fieldErrors)
		httpErrors.ValidationError(w, fieldErrors)
		return
	}

	book := mapBookMetadataToBook(id, metadata)

	if err := h.repo.UpdateBook(r.Context(), book); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to update book: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) UpdateBookImage(w http.ResponseWriter, r *http.Request) {
	id, err := requestx.ParsePathID(r, "id")
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	file, header, fieldErrors := validateBookImage(r)
	if len(fieldErrors) > 0 {
		httpErrors.ValidationError(w, fieldErrors)
		return
	}
	defer file.Close()

	warning, err := h.repo.UpdateBookImage(r.Context(), id, file, header)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to update book image: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}
	if warning != "" {
		h.log.Warn("%s", warning)
	}

	h.log.Info("updated image to %s", header.Filename)

	book, err := h.repo.GetBookByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to fetch book after image update: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := requestx.ParsePathID(r, "id")
	if err != nil {
		httpErrors.NotFound(w)
		return
	}

	warning, err := h.repo.DeleteBook(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpErrors.NotFound(w)
			return
		}
		h.log.Error("failed to delete book: %v", err)
		httpErrors.ServerError(w, http.StatusInternalServerError)
		return
	}
	if warning != "" {
		h.log.Warn("%s", warning)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) GetBookImage(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	image := params.ByName("image")

	if image == "" || image != filepath.Base(image) || strings.Contains(image, "..") {
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	coversDir := filepath.Join(imageDir(), "covers")
	target := filepath.Join(coversDir, image)
	cleanCovers := filepath.Clean(coversDir) + string(os.PathSeparator)
	cleanTarget := filepath.Clean(target)

	if !strings.HasPrefix(cleanTarget, cleanCovers) {
		httpErrors.ClientError(w, http.StatusBadRequest)
		return
	}

	info, err := os.Stat(cleanTarget)
	if err != nil || !info.Mode().IsRegular() {
		httpErrors.NotFound(w)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, cleanTarget)
}

func imageDir() string {
	if dir := os.Getenv("IMAGE_DIR"); dir != "" {
		return dir
	}
	return defaultImageDir
}

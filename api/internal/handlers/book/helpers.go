package book

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/requestx"
	"booksonhooks.ca/internal/validator"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func (h *BookHandler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to write json: %v", err)
	}
}

func mapBookUpdateRequestToBook(id int64, req domain.BookUpdateRequest) *domain.Book {
	return &domain.Book{
		ID:      id,
		Title:   req.Title,
		Author:  req.Author,
		Summary: req.Summary,
		Price:   req.Price,
	}
}

func mapBookUpdateRequestToMetadata(req domain.BookUpdateRequest) bookCreateRequest {
	return bookCreateRequest{
		Title:   requestx.NormalizeText(req.Title),
		Author:  requestx.NormalizeText(req.Author),
		Summary: requestx.NormalizeText(req.Summary),
		Price:   requestx.NormalizeText(req.Price),
	}
}

func mapBookMetadataToBook(id int64, req bookCreateRequest) *domain.Book {
	return &domain.Book{
		ID:      id,
		Title:   req.Title,
		Author:  req.Author,
		Summary: req.Summary,
		Price:   req.Price,
	}
}

type bookCreateRequest struct {
	Title   string
	Author  string
	Summary string
	Price   string
}

func validateCreateBookMetadata(r *http.Request) (bookCreateRequest, map[string]string) {
	req := bookCreateRequest{
		Title:   requestx.NormalizeText(r.PostFormValue("title")),
		Author:  requestx.NormalizeText(r.PostFormValue("author")),
		Summary: requestx.NormalizeText(r.PostFormValue("summary")),
		Price:   requestx.NormalizeText(r.PostFormValue("price")),
	}

	return req, validateBookMetadata(req)
}

func validateBookMetadata(req bookCreateRequest) map[string]string {
	fieldErrors := make(map[string]string)

	if !validator.NotBlank(req.Title) {
		fieldErrors["title"] = "This field cannot be blank"
	} else if !validator.MaxChars(req.Title, 100) {
		fieldErrors["title"] = "This field can have a maximum of 100 characters"
	}

	if !validator.NotBlank(req.Author) {
		fieldErrors["author"] = "This field cannot be blank"
	} else if !validator.MaxChars(req.Author, 100) {
		fieldErrors["author"] = "This field can have a maximum of 100 characters"
	}

	if !validator.NotBlank(req.Summary) {
		fieldErrors["summary"] = "This field cannot be blank"
	}

	if !validator.ValidDollarValue(req.Price) {
		fieldErrors["price"] = "Must be a valid dollar value. IE 10.99 or $10.99"
	}

	return fieldErrors
}

func validateBookImage(r *http.Request) (multipart.File, *multipart.FileHeader, map[string]string) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, nil, map[string]string{
			"image": "Failed to upload file",
		}
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		_ = file.Close()
		return nil, nil, map[string]string{
			"image": "Image must be a .jpg, .jpeg, .png, or .webp file",
		}
	}

	return file, header, nil
}

func mapBookCreateRequestToBook(req bookCreateRequest) *domain.Book {
	return &domain.Book{
		Title:   req.Title,
		Author:  req.Author,
		Summary: req.Summary,
		Price:   req.Price,
	}
}

func (h *BookHandler) recordBookMetric(r *http.Request, bookID int64) {
	var (
		isQr      bool
		machineID int64
		source    string
		sessionID string
		hasErrors bool
	)

	parsedQr, err := requestx.ParseOptionalQueryBool(r, "is_qr")
	if err != nil {
		h.log.Warn("skipping book metric: invalid is_qr book_id=%d raw=%q err=%v",
			bookID, requestx.QueryString(r, "is_qr"), err)
		hasErrors = true
	} else {
		isQr = parsedQr
	}

	parsedMachineID, err := requestx.ParseOptionalQueryInt64(r, "machine")
	if err != nil {
		h.log.Warn("skipping book metric: invalid machine book_id=%d raw=%q err=%v",
			bookID, requestx.QueryString(r, "machine"), err)
		hasErrors = true
	} else {
		machineID = parsedMachineID
	}

	parsedSource, err := requestx.ParseOptionalQuerySource(r, "source")
	if err != nil {
		h.log.Warn("skipping book metric: invalid source book_id=%d raw=%q err=%v",
			bookID, requestx.QueryString(r, "source"), err)
		hasErrors = true
	} else {
		source = parsedSource
	}

	parsedSessionID, err := requestx.ParseOptionalQuerySessionID(r, "session_id")
	if err != nil {
		h.log.Warn("skipping book metric: invalid session_id book_id=%d raw=%q err=%v",
			bookID, requestx.QueryString(r, "session_id"), err)
		hasErrors = true
	} else {
		sessionID = parsedSessionID
	}

	if hasErrors || machineID <= 0 {
		return
	}

	if _, err := h.repo.InsertBookMetric(r.Context(), bookID, machineID, isQr, source, sessionID); err != nil {
		h.log.Error("failed to insert book metric: book_id=%d machine_id=%d is_qr=%t source=%q session_id=%q err=%v",
			bookID, machineID, isQr, source, sessionID, err)
	}
}

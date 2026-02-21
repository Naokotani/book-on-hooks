package book

import (
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/requestx"
	"booksonhooks.ca/internal/validator"
	"encoding/json"
	"mime/multipart"
	"net/http"
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

type bookCreateRequest struct {
	Title   string
	Author  string
	Summary string
	Price   string
}

func validateCreateBookMetadata(r *http.Request) (bookCreateRequest, map[string]string) {
	req := bookCreateRequest{
		Title:   strings.TrimSpace(r.PostFormValue("title")),
		Author:  strings.TrimSpace(r.PostFormValue("author")),
		Summary: strings.TrimSpace(r.PostFormValue("summary")),
		Price:   strings.TrimSpace(r.PostFormValue("price")),
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
		hasErrors bool
	)

	parsedQr, err := requestx.ParseOptionalBool(r.URL.Query().Get("is_qr"))
	if err != nil {
		h.log.Warn("skipping book metric: invalid is_qr book_id=%d raw=%q err=%v",
			bookID, r.URL.Query().Get("is_qr"), err)
		hasErrors = true
	} else {
		isQr = parsedQr
	}

	parsedMachineID, err := requestx.ParseOptionalInt64(r.URL.Query().Get("machine"))
	if err != nil {
		h.log.Warn("skipping book metric: invalid machine book_id=%d raw=%q err=%v",
			bookID, r.URL.Query().Get("machine"), err)
		hasErrors = true
	} else {
		machineID = parsedMachineID
	}

	parsedSource, err := requestx.ParseOptionalSource(r.URL.Query().Get("source"))
	if err != nil {
		h.log.Warn("skipping book metric: invalid source book_id=%d raw=%q err=%v",
			bookID, r.URL.Query().Get("source"), err)
		hasErrors = true
	} else {
		source = parsedSource
	}

	if hasErrors || machineID <= 0 {
		return
	}

	if _, err := h.repo.InsertBookMetric(r.Context(), bookID, machineID, isQr, source); err != nil {
		h.log.Error("failed to insert book metric: book_id=%d machine_id=%d is_qr=%t source=%q err=%v",
			bookID, machineID, isQr, source, err)
	}
}

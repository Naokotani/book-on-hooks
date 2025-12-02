package forms

import (
	"fmt"
	"net/http"
	"strings"

	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/validator"
)

type BookCreateForm struct {
	Title               string `form:"title"`
	Author              string `form:"author"`
	Price               string `form:"price"`
	Summary             string `form:"summary"`
	validator.Validator `form:"-"`
}

func ParseBookForm(f *Form, r *http.Request) (int64, *BookCreateForm, *httpErrors.HTTPError) {
	var form BookCreateForm

	err := f.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to decode form: %w", err)}
	}

	form.Title = strings.TrimSpace(form.Title)
	form.Author = strings.TrimSpace(form.Author)
	form.Price = strings.TrimSpace(form.Price)
	form.Summary = strings.TrimSpace(form.Summary)

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field can have a maximum of 100 characters")
	form.CheckField(validator.NotBlank(form.Author), "author", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Author, 100), "author", "This field can have a maximum of 100 characters")
	form.CheckField(validator.NotBlank(form.Summary), "summar", "This field cannot be blank")
	form.CheckField(validator.ValidDollarValue(form.Price), "price", "Must be a valid dollar value. IE 10.99 or $10.99")

	book := repository.Book{Title: form.Title, Author: form.Author, Summary: form.Summary}

	if !form.Valid() {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to validate form. Book: %v", book)}
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to upload file: %w", err)}
	}
	defer file.Close()

	id, err := f.repo.InsertBook(&book, file, header)

	if err != nil {
		form.CheckField(false, "image", "Failed to upload file")
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusInternalServerError, Err: fmt.Errorf("failed to insert book: %w", err)}
	}
	f.log.Info("Book input: %v\n", book)
	return id, &form, nil
}

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

func BookFormService(f *Form, r *http.Request) (int64, *BookCreateForm, *httpErrors.HTTPError) {
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
	form.CheckField(validator.NotBlank(form.Summary), "summary", "This field cannot be blank")
	form.CheckField(validator.ValidDollarValue(form.Price), "price", "Must be a valid dollar value. IE 10.99 or $10.99")

	book := repository.Book{Title: form.Title, Author: form.Author, Summary: form.Summary, Price: form.Price}

	if !form.Valid() {
		f.log.Warn("Invalid form submitted: %v", form)
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to validate form. Book: %v", book)}
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		f.log.Warn("Failed to get form image for: %v", form)
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to upload file: %w", err)}
	}
	defer file.Close()

	id, err := f.repo.InsertBook(&book, file, header)

	if err != nil {
		f.log.Error("Failed to get form image for: %v", form)
		form.CheckField(false, "image", "Failed to upload file")
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusInternalServerError, Err: fmt.Errorf("failed to insert book: %w", err)}
	}
	return id, &form, nil
}

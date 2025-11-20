package main

import (
	"fmt"
	"net/http"
	"strings"

	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/validator"
)

type bookCreateForm struct {
	Title               string `form:"title"`
	Author              string `form:"author"`
	Price               string `form:"price"`
	Summary             string `form:"summary"`
	Machine             int    `form:"machine"`
	validator.Validator `form:"-"`
}

func parseBookForm(app *application, r *http.Request) (int64, *bookCreateForm, *httpErrors.HTTPError) {

	var form bookCreateForm

	err := app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to decode form: %w", err)}
	}

	form.Title = strings.TrimSpace(form.Title)
	form.Author = strings.TrimSpace(form.Author)
	form.Price = strings.TrimSpace(form.Price)

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field can have a maximum of 100 characters")
	form.CheckField(validator.NotBlank(form.Author), "author", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Author, 100), "author", "This field can have a maximum of 100 characters")
	form.CheckField(validator.ValidDollarValue(form.Price), "price", "Must be a valid dollar value. IE 10.99 or $10.99")

	book := repository.Book{Title: form.Title, Author: form.Author, Price: form.Price}

	if !form.Valid() {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to validate form. Book: %v", book)}
	}

	file, header, err := r.FormFile("cover")
	if err != nil {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to upload file: %w", err)}
	}
	defer file.Close()

	id, err := app.db.InsertBook(&book, file, header)

	if err != nil {
		form.CheckField(false, "cover", "Failed to upload file")
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusInternalServerError, Err: fmt.Errorf("failed to insert book: %w", err)}
	}
	return id, &form, nil
}

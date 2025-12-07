package forms

import (
	"fmt"
	"net/http"
	"strings"

	"booksonhooks.ca/internal/httpErrors"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/validator"
)

type MachineCreateForm struct {
	Location            string `form:"location"`
	Rows                int    `form:"rows"`
	Cols                int    `form:"cols"`
	validator.Validator `form:"-"`
}

func MachineFormService(f *Form, r *http.Request) (int64, *MachineCreateForm, *httpErrors.HTTPError) {
	var form MachineCreateForm

	r.ParseForm()
	err := f.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to decode form: %w", err)}
	}
	f.log.Info("Inserting machine. %v\n", form.Location)
	f.log.Info("post form. %v\n", r.PostForm)

	form.Location = strings.TrimSpace(form.Location)

	form.CheckField(validator.NotBlank(form.Location), "location", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Location, 100), "location", "This field can have a maximum of 100 characters")

	machine := repository.Machine{Location: form.Location}

	if !form.Valid() {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusUnprocessableEntity, Err: fmt.Errorf("failed to validate form. Book: %v", machine)}
	}

	_, err = f.repo.InsertMachine(r.Context(), &machine)

	if err != nil {
		return 0, &form, &httpErrors.HTTPError{Status: http.StatusInternalServerError, Err: fmt.Errorf("failed to insert machine: %w", err)}
	}
	f.log.Info("Machine inserted succesfully. ID: %d, Location: %s \n", machine.ID, machine.Location)
	return machine.ID, &form, nil
}

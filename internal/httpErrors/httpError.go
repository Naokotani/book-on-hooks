package httpErrors

import (
	"net/http"
)

type HTTPError struct {
	Status int
	Err    error
}

func (e *HTTPError) Error() string { return e.Err.Error() }
func (e *HTTPError) UnWrap() error { return e.Err }

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}

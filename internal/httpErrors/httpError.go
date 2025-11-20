package httpErrors

type HTTPError struct {
	Status int
	Err    error
}

func (e *HTTPError) Error() string { return e.Err.Error() }
func (e *HTTPError) UnWrap() error { return e.Err }

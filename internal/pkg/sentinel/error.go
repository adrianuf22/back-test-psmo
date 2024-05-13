package sentinel

import (
	"net/http"
)

var (
	ErrSafeFatal  = []byte("{\"error\":\"internal error\"}")
	ErrUnauth     = &httpError{status: http.StatusUnauthorized, Msg: "invalid token"}
	ErrNotFound   = &httpError{status: http.StatusNotFound, Msg: "not found"}
	ErrDuplicate  = &httpError{status: http.StatusConflict, Msg: "duplicate resource"}
	ErrBadRequest = &httpError{status: http.StatusBadRequest, Msg: "bad input data"}
	ErrInternal   = &httpError{status: http.StatusInternalServerError, Msg: "internal error"}
)

type Error interface {
	HttpError() (int, httpError)
}

type httpError struct {
	status int         `json:"-"`
	Msg    string      `json:"error"`
	Values interface{} `json:"values,omitempty"`
}

func (e httpError) Error() string {
	return e.Msg
}

func (e httpError) HttpError() (int, httpError) {
	return e.status, e
}

func (e *httpError) WithValues(values interface{}) *httpError {
	e.Values = values

	return e
}

type wrappedError struct {
	error
	sentinel *httpError
}

func (e wrappedError) Is(err error) bool {
	return e.sentinel == err
}

func (e wrappedError) HttpError() (int, httpError) {
	return e.sentinel.HttpError()
}

func WrapError(err error, sentinel *httpError) error {
	return wrappedError{error: err, sentinel: sentinel}
}

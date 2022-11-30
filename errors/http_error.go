package errors

import (
	"errors"
	"net/http"
)

type HttpError struct {
	status int
	err    error
}

func (he *HttpError) Error() string {
	return he.err.Error()
}

func (he *HttpError) Status() int {
	return he.status
}

func New(errMsg string, status int) error {
	return &HttpError{status, errors.New(errMsg)}
}

var BadRequest = &HttpError{
	http.StatusBadRequest,
	errors.New("bad request"),
}

var NotFound = &HttpError{
	http.StatusNotFound,
	errors.New("resource not found")}

var Conflict = &HttpError{
	http.StatusConflict,
	errors.New("resource already exists")}

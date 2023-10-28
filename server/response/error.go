package response

import (
	"fmt"
	"net/http"
)

// swagger:model
type Error struct {
	StatusCode int
	Message    map[string][]string
	Err        error
}

func New(err error, status int, msg map[string][]string) Error {
	return Error{
		StatusCode: status,
		Err:        err,
		Message:    msg,
	}
}

func (r Error) Error() string {
	return r.Err.Error()
}

func (r Error) Unwrap() error {
	return r.Err
}

func BadRequest(err error) Error {
	return New(
		err,
		http.StatusBadRequest,
		BadRequestMsg,
	)
}

func BadRequestMsg(msg string) Error {
	return New(
		fmt.Errorf(""),
		http.StatusBadRequest,
		map[string][]string{"__error__": {msg}},
	)
}

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

// BadRequest shortcut to return http.StatusBadRequest with custom error and msg
func BadRequest(err error, msg string) Error {
	if err == nil {
		err = fmt.Errorf("")
	}
	finalMsg := map[string][]string{"__error__": {msg}}
	if msg == "" {
		finalMsg = MsgBadRequest
	}
	return New(
		err,
		http.StatusBadRequest,
		finalMsg,
	)
}

// NotFound shortcut to return http.StatusNotFound with custom error and msg
func NotFound(err error, msg string) Error {
	if err == nil {
		err = fmt.Errorf("")
	}
	finalMsg := map[string][]string{"__error__": {msg}}
	if msg == "" {
		finalMsg = MsgNotFound
	}
	return New(
		err,
		http.StatusNotFound,
		finalMsg,
	)
}

// NotFound shortcut to return http.StatusNotFound with custom error and msg
func InternalError(err error, msg string) Error {
	if err == nil {
		err = fmt.Errorf("")
	}
	finalMsg := map[string][]string{"__error__": {msg}}
	if msg == "" {
		finalMsg = MsgInternalErr
	}
	return New(
		err,
		http.StatusInternalServerError,
		finalMsg,
	)
}

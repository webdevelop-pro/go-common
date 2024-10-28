package response

import (
	"net/http"

	"github.com/pkg/errors"
)

// BadRequest shortcut to return http.StatusBadRequest with custom error and msg
func BadRequest(err error, msg string) Error {
	if err == nil {
		err = errors.New("")
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
		err = errors.New("")
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

// InternalError shortcut to return http.StatusInternalServerError with custom error and msg
func InternalError(err error, msg string) Error {
	if err == nil {
		err = errors.New("")
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

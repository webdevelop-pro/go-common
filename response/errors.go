package response

import (
	"fmt"
	"net/http"
)

var (
	DefaultErrBadRequest = NewError(
		fmt.Errorf("bad request"), http.StatusBadRequest, MsgBadRequest,
	)
	DefaultErrInternalError = NewError(
		fmt.Errorf("internal error"), http.StatusInternalServerError, MsgInternalErr,
	)
)

func ErrBadRequest(err error) *Error {
	return NewError(err, http.StatusBadRequest, map[string][]string{"__error__": {err.Error()}})
}

func ErrInternalError(err error) *Error {
	return NewError(err, http.StatusInternalServerError, map[string][]string{"__error__": {err.Error()}})
}

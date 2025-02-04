package response

import (
	"fmt"
)

var (
	DefaultErrBadRequest = NewError(
		fmt.Errorf("bad request"), StatusBadRequest, MsgBadRequest,
	)
	DefaultErrUnauthorized = NewError(
		fmt.Errorf("unauthorized"), StatusUnauthorized, MsgUnauthorized,
	)
	DefaultErrNotFound = NewError(
		fmt.Errorf("not found"), StatusNotFound, MsgNotFound,
	)
	DefaultErrInternalError = NewError(
		fmt.Errorf("internal error"), StatusInternalError, MsgInternalErr,
	)
)

func ErrBadRequest(err error) *Error {
	return NewError(err, StatusBadRequest, map[string][]string{"__error__": {err.Error()}})
}

func ErrUnauthorized(err error) *Error {
	return NewError(err, StatusUnauthorized, map[string][]string{"__error__": {err.Error()}})
}

func ErrInternalError(err error) *Error {
	return NewError(err, StatusInternalError, map[string][]string{"__error__": {err.Error()}})
}

package response

import (
	"fmt"
)

var (
	DefaultErrBadRequest = NewError(
		fmt.Errorf("bad request"), StatusBadRequest, MsgBadRequest,
	)
	DefaultErrInternalError = NewError(
		fmt.Errorf("internal error"), StatusInternalError, MsgInternalErr,
	)

	DefaultErrNotFound = NewError(
		fmt.Errorf("not found"), StatusNotFound, MsgNotFound,
	)
)

func ErrBadRequest(err error) *Error {
	return NewError(err, StatusBadRequest, map[string][]string{"__error__": {err.Error()}})
}

func ErrInternalError(err error) *Error {
	return NewError(err, StatusInternalError, map[string][]string{"__error__": {err.Error()}})
}

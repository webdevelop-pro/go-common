package response

import "net/http"

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

func BadRequest(err error) Error {
	return New(
		err,
		http.StatusBadRequest,
		BadRequestMsg,
	)
}

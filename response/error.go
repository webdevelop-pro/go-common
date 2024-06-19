package response

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

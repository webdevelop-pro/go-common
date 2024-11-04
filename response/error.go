package response

// swagger:model
type Error struct {
	StatusCode int   `json:"-"`
	Message    any   `json:"message"`
	Err        error `json:"-"`
}

func NewError(err error, args ...any) *Error {
	newError := Error{
		Err: err,
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			newError.StatusCode = v
		case string:
			newError.Message = v
		case map[string][]string:
			newError.Message = v
		}
	}

	return &newError
}

func (r *Error) Error() string {
	return r.Err.Error()
}

func (r *Error) Unwrap() error {
	return r.Err
}

func (r *Error) GetMessageFromMap(key string) ([]string, bool) {
	_, ok := r.Message.(map[string][]string)
	if !ok {
		return nil, false
	}

	_, ok = r.Message.(map[string][]string)[key]
	if !ok {
		return nil, false
	}

	return r.Message.(map[string][]string)[key], true
}

func (r *Error) AddMessageToMap(key string, value string) {
	if r.Message == nil {
		r.Message = make(map[string][]string)
	}

	_, ok := r.Message.(map[string][]string)[key]
	if !ok {
		r.Message.(map[string][]string)[key] = []string{}
	}

	r.Message.(map[string][]string)[key] = append(r.Message.(map[string][]string)[key], value)
}

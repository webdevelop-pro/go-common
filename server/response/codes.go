package response

var (
	BadRequestMsg = map[string][]string{"__error__": {"bad request"}}
	// BadAuth signalizes lack of token or other authorization data
	NotAuthorizedMsg = map[string][]string{"__error__": {"authorization error"}}
	// InternalError server side error
	InternalErrMsg = map[string][]string{"__error__": {"internal/server error"}}
)

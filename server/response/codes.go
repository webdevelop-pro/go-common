package response

var (
	// MsgBadRequest used to indicate error in incoming data
	MsgBadRequest = map[string][]string{"__error__": {"bad request"}}
	// MsgNotFound typically used when element haven't been found
	MsgNotFound = map[string][]string{"__error__": {"not found"}}
	// MsgUnauthorized signalizes lack of token or other authorization data
	MsgUnauthorized = map[string][]string{"__error__": {"authorization error"}}
	// MsgForbidden used when user does not have any permissions to perform action
	MsgForbidden = map[string][]string{"__error__": {"forbidden error"}}
	// MsgInternalErr server side error
	MsgInternalErr = map[string][]string{"__error__": {"internal/server error"}}
)

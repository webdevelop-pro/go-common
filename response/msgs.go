//nolint:gochecknoglobals
package response

var (
	// MsgBadRequest used to indicate error in incoming data
	MsgBadRequest = map[string][]string{"__error__": {"Bad request."}}
	// MsgNotFound typically used when element haven't been found
	MsgNotFound = map[string][]string{"__error__": {"Not found."}}
	// MsgUnauthorized signalizes lack of token or other authorization data
	MsgUnauthorized = map[string][]string{"__error__": {"Authorization error."}}
	// MsgForbidden used when user does not have any permissions to perform action
	MsgForbidden = map[string][]string{"__error__": {"Forbidden error."}}
	// MsgInternalErr server side error
	MsgInternalErr = map[string][]string{"__error__": {"Internal/server error."}}
)

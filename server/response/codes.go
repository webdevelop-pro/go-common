package response

const (
	// BadAuth signalizes lack of token or other authorization data
	BadAuth = "unauthorized"
	// InternalError server side error
	InternalError = "internal/server"
)

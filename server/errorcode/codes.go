package errorcode

const (
	// BadAuth signalizes lack of token or other authorization data
	BadAuth = "user/unauthorized"
	// NoEntity means that there is no such entity
	NoEntity = "user/no-entity"
	// InternalError server side error
	InternalError = "internal/server"
)

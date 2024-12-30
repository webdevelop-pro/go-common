package response

import "net/http"

const (
	StatusBadRequest    = http.StatusBadRequest
	StatusUnauthorized  = http.StatusUnauthorized
	StatusNotFound      = http.StatusNotFound
	StatusInternalError = http.StatusInternalServerError
)

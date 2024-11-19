package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// DefaultSkipper returns false which processes the middleware.
func FileAndHealtchCheckSkipper(c echo.Context) bool {
	// skip healthcheck requests
	// ToDo
	// CreateConfig to ignore log URLS
	if c.Request().URL.Path == "/healthcheck" || c.Request().URL.Path == "/metrics" {
		return true
	}
	// do not dump body for multipart requests
	contentType := c.Request().Header.Get("Content-Type")
	if len(contentType) > 10 && contentType[0:10] == "multipart/" {
		return true
	}
	return false
}

// Writes down every request to the log file
// Useful for debagging
func BodyDumpHandler(c echo.Context, incoming []byte, outcoming []byte) {
	log := zerolog.Ctx(c.Request().Context())
	log.Trace().Str("path", c.Request().RequestURI).
		Interface("headers", c.Request().Header).
		Interface("body", string(incoming)).Msg("incoming request")
	// ToDo
	// headers = Fix error marshaling error: json: unsupported type: func() http.Header
	log.Trace().Interface("headers", c.Response().Header).
		Interface("body", string(outcoming)).Msg("outcoming request")
}

package middleware

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	logger "github.com/webdevelop-pro/go-logger"
)

// Writes down every request to the log file
// Useful for debagging
func LogRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ignore healthcheck requests
		if c.Request().URL.Path != "/healthcheck" {
			// create sub logger
			log := logger.NewComponentLogger("log-requests", c)
			// enrich context
			raw, _ := io.ReadAll(c.Request().Body)
			c.Request().Body = io.NopCloser(bytes.NewReader(raw))
			log.Trace().Str("path", c.Request().RequestURI).
				Interface("headers", c.Request().Header).
				Interface("body", string(raw)).Msg("raw request")
		}
		// next handler
		return next(c)
	}
}

// ToDo
// Ability to set up log level in config file

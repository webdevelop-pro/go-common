package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/logger"
)

// SetLogger adds logger to context
func SetLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get request's context
		ctx := c.Request().Context()
		// create sub logger
		log := logger.NewComponentLogger("echo", c)
		// add logger to context
		ctx = log.WithContext(ctx)
		// enrich context
		c.SetRequest(c.Request().WithContext(ctx))
		// next handler
		return next(c)
	}
}

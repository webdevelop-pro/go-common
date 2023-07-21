// Package to add timestamp with milliseconds to track latency
package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// ContextKey is the key used to lookup incoming time
	// request from the echo.Context.
	RequestTimeContextKey = "request-created-at"
)

// SetRequestTime sets an initial request time
func SetRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentTime := time.Now()
		c.Set(RequestTimeContextKey, currentTime)
		return next(c)
	}
}

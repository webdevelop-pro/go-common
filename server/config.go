package server

import (
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"github.com/webdevelop-pro/go-common/server/middleware"
)

// Config is struct to configure HTTP server
type Config struct {
	Host string `required:"true"`
	Port string `required:"true"`
}

func DefaultMiddlewares(middlewares ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	defaults := []echo.MiddlewareFunc{
		// Set context logger
		middleware.SetIPAddress,
		middleware.DefaultCTXValues,
		middleware.SetRequestTime,
		middleware.LogRequests,
		// Trace ID middleware generates a unique id for a request.
		echoMW.RequestIDWithConfig(echoMW.RequestIDConfig{
			RequestIDHandler: func(c echo.Context, requestID string) {
				c.Set(echo.HeaderXRequestID, requestID)
			},
		}),
	}
	if len(middlewares) > 0 {
		defaults = append(defaults, middlewares...)
	}
	return defaults
}

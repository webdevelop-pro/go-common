package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/webdevelop-pro/lib/constants"
)

// SetTraceID adds trace_id to context and X-TRACE-ID to headers
func SetTraceID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// generate trace id
		trcID := uuid.New().String()
		// get request's context
		ctx := c.Request().Context()
		// add random trace id
		ctx = context.WithValue(ctx, constants.CtxTraceID, trcID)
		// create sub logger
		l := log.Ctx(ctx).With().Str("trace_id", trcID).Logger()
		// add logger to context
		ctx = l.WithContext(ctx)
		// enrich context
		c.SetRequest(c.Request().WithContext(ctx))
		// add as header
		c.Response().Header().Set("X-TRACE-ID", trcID)
		// next handler
		return next(c)
	}
}

// GetTraceID is a function which  extract trace id from context
func GetTraceID(ctx context.Context) string {
	trcID, ok := ctx.Value(constants.CtxTraceID).(string)
	if !ok {
		return ""
	}
	return trcID
}

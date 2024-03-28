package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
)

func SetDefaultCTX(echoCtx echo.Context) context.Context {
	ctx := SetDefaultHTTPCtx(echoCtx.Request().Context(), echoCtx.Request().Header)
	if ctx.Value(keys.RequestID).(string) == "" {
		requestID := echoCtx.Get(echo.HeaderXRequestID).(string)
		ctx = keys.SetCtxValue(ctx, keys.RequestID, requestID)
	}
	return ctx
}

func SetDefaultHTTPCtx(ctx context.Context, headers http.Header) context.Context {
	requestID := headers.Get(echo.HeaderXRequestID)
	IP := GetIpAddress(headers)

	ctx = keys.SetCtxValue(ctx, keys.RequestID, requestID)
	ctx = keys.SetCtxValue(ctx, keys.IPAddress, IP)
	return ctx
}

// Set values in ctx for
// RequestID, IPAddress
func DefaultCTXValues(next echo.HandlerFunc) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		ctx := SetDefaultCTX(echoCtx)
		echoCtx.SetRequest(echoCtx.Request().WithContext(ctx))

		// next handler
		return next(echoCtx)
	}
}

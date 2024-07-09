package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
)

func SetDefaultCTX(echoCtx echo.Context) context.Context {
	ctx := keys.SetDefaultHTTPCtx(echoCtx.Request().Context(), echoCtx.Request().Header)
	if val, ok := ctx.Value(keys.RequestID).(string); ok && val == "" {
		if requestID, ok := echoCtx.Get(echo.HeaderXRequestID).(string); ok {
			ctx = keys.SetCtxValue(ctx, keys.RequestID, requestID)
		}
	}
	return ctx
}

package keys

import (
	"context"
)

type (
	ContextKey rune
	ContextStr string
)

const (
	CtxTraceID ContextKey = iota
	RequestID
	IPAddress
	MSGID
	IdentityID
	LogInfo
	RequestLogID

	RequestIDStr ContextStr = "X-Request-Id"
	IPAddressStr ContextStr = "IP-Address"

	XOriginalForwardedFor = "X-Original-Forwarded-For"
	XForwardedFor         = "X-Forwarded-For"
	XRealIP               = "X-Real-IP"
)

func GetAsString(ctx context.Context, key any) string {
	val, ok := ctx.Value(key).(string)
	if ok {
		return val
	}
	return ""
}

func GetCtxValue(ctx context.Context, key any) any {
	return ctx.Value(key)
}

func SetCtxValue(ctx context.Context, key ContextKey, value any) context.Context {
	ctx = context.WithValue(ctx, key, value)

	return ctx
}

func SetCtxValues(ctx context.Context, values map[ContextKey]any) context.Context {
	for key, value := range values {
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx
}

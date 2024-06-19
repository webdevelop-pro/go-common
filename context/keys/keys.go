package keys

import (
	"context"
)

type ContextKey rune

const (
	CtxTraceID ContextKey = iota
	RequestID
	IPAddress
	MSGID
	IdentityID
	LogInfo
	RequestLogID
	LogObjectType
	LogObjectID
)

func GetAsString(ctx context.Context, key ContextKey) string {
	val, ok := ctx.Value(key).(string)
	if ok {
		return val
	}
	return ""
}

func GetCtxValue(ctx context.Context, key ContextKey) any {
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

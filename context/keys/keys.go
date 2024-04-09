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

func GetCtxValue(ctx context.Context, key ContextKey) any {
	return ctx.Value(key)
}

func SetCtxValue(ctx context.Context, key ContextKey, value any) context.Context {
	ctx = context.WithValue(ctx, key, value)

	return ctx
}

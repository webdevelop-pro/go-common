package keys

import (
	"context"
	"net/http"
	"strings"
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

	RequestIDStr = "X-Request-Id"
	IPAddressStr = "IP-Address"
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

func GetIPAddress(headers http.Header) string {
	// ToDo
	// Use echo context.RealIP()
	ip := "127.0.0.1"
	if xOFF := headers.Get("X-Original-Forwarded-For"); xOFF != "" {
		i := strings.Index(xOFF, ", ")
		if i == -1 {
			i = len(xOFF)
		}
		ip = xOFF[:i]
	} else if xFF := headers.Get("X-Forwarded-For"); xFF != "" {
		i := strings.Index(xFF, ", ")
		if i == -1 {
			i = len(xFF)
		}
		ip = xFF[:i]
	} else if xrIP := headers.Get("X-Real-IP"); xrIP != "" {
		ip = xrIP
	}
	return ip
}

// Set values in ctx for
// RequestID, IPAddress
func SetDefaultHTTPCtx(ctx context.Context, headers http.Header) context.Context {
	// so we don't need echo here
	// requestID := headers.Get(echo.HeaderXRequestID)
	requestID := headers.Get("X-Request-Id")
	IP := GetIPAddress(headers)

	ctx = SetCtxValue(ctx, RequestID, requestID)
	ctx = SetCtxValue(ctx, IPAddress, IP)
	return ctx
}

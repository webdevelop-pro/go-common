package httputils

import (
	"context"
	"net/http"
	"strings"
)

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

	ctx = SetCtxValue(ctx, KeyRequestID, requestID)
	ctx = SetCtxValue(ctx, KeyIPAddress, IP)
	return ctx
}

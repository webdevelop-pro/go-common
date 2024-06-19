package pclient

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/verser"
)

func GetIPAddress(headers http.Header) string {
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

func SetDefaultEventCtx(ctx context.Context, event Event) context.Context {
	ctx = keys.SetCtxValue(ctx, keys.RequestID, event.RequestID)
	ctx = keys.SetCtxValue(ctx, keys.IPAddress, event.IPAddress)
	ctx = keys.SetCtxValue(ctx, keys.MSGID, event.ID)

	logInfo := logger.ServiceContext{
		Service: verser.GetService(),
		Version: verser.GetVersion(),
		SourceReference: &logger.SourceReference{
			Repository: verser.GetRepository(),
			RevisionID: verser.GetRevisionID(),
		},
		MSGID: event.ID,
	}

	ctx = keys.SetCtxValue(ctx, keys.LogInfo, logInfo)

	return ctx
}

func SetDefaultWebhookCtx(ctx context.Context, webhook Webhook) context.Context {
	headers := http.Header(webhook.Headers)

	requestID := headers.Get(echo.HeaderXRequestID)
	IP := GetIPAddress(headers)

	ctx = keys.SetCtxValue(ctx, keys.RequestID, requestID)
	ctx = keys.SetCtxValue(ctx, keys.IPAddress, IP)
	ctx = keys.SetCtxValue(ctx, keys.MSGID, webhook.ID)

	logInfo := logger.ServiceContext{
		Service: verser.GetService(),
		Version: verser.GetVersion(),
		SourceReference: &logger.SourceReference{
			Repository: verser.GetRepository(),
			RevisionID: verser.GetRevisionID(),
		},
		MSGID: webhook.ID,
	}

	ctx = keys.SetCtxValue(ctx, keys.LogInfo, logInfo)

	return ctx
}

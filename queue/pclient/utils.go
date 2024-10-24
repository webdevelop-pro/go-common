package pclient

import (
	"context"
	"net/http"

	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/httputils"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/verser"
)

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

	requestID := headers.Get(string(keys.RequestIDStr))
	IP := httputils.GetIPAddress(headers)

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

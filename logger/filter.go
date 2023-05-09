package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/verser"
)

const (
	errorTypeKey   = "@type"
	errorTypeValue = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type GoogleCloudHook struct {
	skip       bool
	reqContext echo.Context
}

func (h GoogleCloudHook) echoContextToServiceContext() ServiceContext {
	serviceCtx := ServiceContext{
		HttpRequest: &HttpRequestContext{
			Method:    h.reqContext.Request().Method,
			URL:       h.reqContext.Request().URL.String(),
			UserAgent: h.reqContext.Request().UserAgent(),
			Referrer:  h.reqContext.Request().Referer(),
			RemoteIp:  h.reqContext.Request().RemoteAddr,
		},
	}
	if service := verser.GetService(); service != "" {
		serviceCtx.Service = service
	}

	if version := verser.GetVersion(); version != "" {
		serviceCtx.Version = version
	}

	repository := verser.GetRepository()
	revisionID := verser.GetRevisionID()
	if repository != "" || revisionID != "" {
		serviceCtx.SourceReference = &SourceReference{
			Repository: repository,
			RevisionID: revisionID,
		}
	}

	if h.reqContext != nil {
		if user := h.reqContext.Get("user"); user != nil {
			serviceCtx.User = user.(string)
		}
	}
	return serviceCtx
}

func (h GoogleCloudHook) Run(e *zerolog.Event, level zerolog.Level, s string) {
	if h.skip {
		return
	}

	switch level {
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		if h.reqContext != nil {
			e.Interface("serviceContext", h.echoContextToServiceContext())
		}
		e.Str(errorTypeKey, errorTypeValue)
	}
}

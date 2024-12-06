package logger

import (
	"context"

	"github.com/rs/zerolog"
)

const (
	ServiceContextInfo = "service_context_info"
)

// Logger is wrapper struct around logger.Logger that adds some custom functionality
type Logger struct {
	zerolog.Logger
}

func (l *Logger) Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// ServiceContext contain info for all logs
type ServiceContext struct {
	Service         string              `json:"service"`
	Version         string              `json:"version"`
	User            string              `json:"user,omitempty"`
	RequestID       string              `json:"request_id,omitempty"`
	MSGID           string              `json:"msg_id,omitempty"`
	HTTPRequest     *HTTPRequestContext `json:"httpRequest,omitempty"`
	SourceReference *SourceReference    `json:"sourceReference,omitempty"`
}

// SourceReference repository name and revision id
type SourceReference struct {
	Repository string `json:"repository"`
	RevisionID string `json:"revisionId"`
}

// HTTPRequestContext http request context
type HTTPRequestContext struct {
	Method             string `json:"method"`
	URL                string `json:"url"`
	UserAgent          string `json:"userAgent"`
	Referrer           string `json:"referrer"`
	ResponseStatusCode int    `json:"responseStatusCode"`
	RemoteIP           string `json:"remoteIp"`
}

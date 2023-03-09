package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/lib/configurator"
	"github.com/webdevelop-pro/lib/verser"
)

// Logger is wrapper struct around logger.Logger that adds some custom functionality
type Logger struct {
	zerolog.Logger
}

// ServiceContext contain info for all logs
type ServiceContext struct {
	Service         string              `json:"service"`
	Version         string              `json:"version"`
	User            string              `json:"user,omitempty"`
	HttpRequest     *HttpRequestContext `json:"httpRequest,omitempty"`
	SourceReference *SourceReference    `json:"sourceReference,omitempty"`
}

// SourceReference repositary name and revision id
type SourceReference struct {
	Repository string `json:"repository"`
	RevisionID string `json:"revisionId"`
}

// HttpRequestContext http request context
type HttpRequestContext struct {
	Method             string `json:"method"`
	URL                string `json:"url"`
	UserAgent          string `json:"userAgent"`
	Referrer           string `json:"referrer"`
	ResponseStatusCode int    `json:"responseStatusCode"`
	RemoteIp           string `json:"remoteIp"`
}

// Printf is implementation of fx.Printer
func (l Logger) Printf(s string, args ...interface{}) {
	l.Info().Msgf(s, args...)
}

// NewLogger return logger instance
func NewLogger(component string, output io.Writer, conf *configurator.Configurator) Logger {
	cfg := conf.New("logger", &Config{}).(*Config)

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else if output == nil {
		output = os.Stdout
	}

	l := zerolog.
		New(output).
		Level(level).
		Hook(SeverityHook{}).
		Hook(TypeHook{skip: cfg.LogConsole}).
		With().Timestamp()

	if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
		l = l.Caller()
	}

	if component != "" {
		l = l.Str("component", component)
	}

	serviceCtx := ServiceContext{
		SourceReference: &SourceReference{},
		HttpRequest:     &HttpRequestContext{},
	}

	if service := verser.GetService(); service != "" {
		serviceCtx.Service = service
		l = l.Str("service", service)
	}

	if version := verser.GetVersion(); version != "" {
		serviceCtx.Version = version
		l = l.Str("version", version)
	}

	if repository := verser.GetRepository(); repository != "" {
		serviceCtx.SourceReference.Repository = repository
		l = l.Str("repository", repository)
	}

	if revisionID := verser.GetRevisionID(); revisionID != "" {
		serviceCtx.SourceReference.RevisionID = revisionID
		l = l.Str("revisionID", revisionID)
	}

	if serviceCtx.Service != "" || serviceCtx.Version != "" {
		l = l.Interface("serviceContext", serviceCtx)
	}

	return Logger{l.Logger()}
}

// NewDefaultLogger return default logger instance
func NewDefaultLogger() Logger {
	return NewLogger("", os.Stdout, configurator.NewConfigurator())
}

// NewDefaultComponentLogger return default logger instance with custom component
func NewDefaultComponentLogger(component string) Logger {
	return NewLogger(component, os.Stdout, configurator.NewConfigurator())
}

// NewDefaultConsoleLogger return default logger instance
func NewDefaultConsoleLogger() Logger {
	return NewLogger("", zerolog.NewConsoleWriter(), configurator.NewConfigurator())
}

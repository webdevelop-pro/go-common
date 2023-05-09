package logger

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/webdevelop-pro/go-common/configurator"
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

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// Printf is implementation of fx.Printer
func (l Logger) Printf(s string, args ...interface{}) {
	l.Info().Msgf(s, args...)
}

// NewLogger return logger instance
func NewLogger(component string, logLevel string, output io.Writer, c echo.Context) Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	l := zerolog.
		New(output).
		Level(level).
		// Hook(SeverityHook{}).
		// Hook(GoogleCloudHook{reqContext: c, skip: false}). // output == ConsoleWriter ?
		With().Stack().Timestamp()

	// if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
	// l = l.Caller()
	// }

	if component != "" {
		l = l.Str("component", component)
	}

	if err != nil {
		ll := l.Logger()
		ll.Error().Err(err).Interface("level", logLevel).Msg("cannot parse log level, using default info")
	}

	return Logger{l.Logger()}
}

// DefaultStdoutLogger return default logger instance
func DefaultStdoutLogger(logLevel string, c echo.Context) Logger {
	return NewLogger("default", logLevel, os.Stdout, c)
}

// NewComponentLogger return default logger instance with custom component
func NewComponentLogger(component string, c echo.Context) Logger {
	conf := configurator.NewConfigurator()
	cfg := conf.New("logger", &Config{}).(*Config)

	var output io.Writer
	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else {
		output = os.Stdout
	}

	return NewLogger(component, cfg.LogLevel, output, c)
}

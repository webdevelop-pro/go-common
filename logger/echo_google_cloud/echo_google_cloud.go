package echogooglecloud

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"
)

const (
	errorTypeKey   = "@type"
	errorTypeValue = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type EchoGoogleCloud struct {
	skip bool
}

func (h EchoGoogleCloud) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if h.skip {
		return
	}

	switch level {
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		e.Str(errorTypeKey, errorTypeValue)
	case zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.NoLevel, zerolog.Disabled, zerolog.TraceLevel:
	}
}

// NewLogger return logger instance
func NewEchoGCLogger(c context.Context, component string, logLevel string, output io.Writer) logger.Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// locally we don't need details for google cloud logging
	_, skipGoogleHook := output.(zerolog.ConsoleWriter)
	l := zerolog.
		New(output).
		Level(level).
		Hook(logger.SeverityHook{}).
		Hook(logger.ContextHook{}).
		Hook(EchoGoogleCloud{skip: !skipGoogleHook}).
		With().Timestamp()

	// if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
	// l = l.Caller()
	// }

	if component != "" {
		l = l.Str("component", component)
	}

	if c != nil {
		l = l.Ctx(c)
	}

	if err != nil {
		ll := l.Logger()
		ll.Error().Err(err).Interface("level", logLevel).Msg("cannot parse log level, using default info")
	}

	return logger.Logger{Logger: l.Logger()}
}

// DefaultStdoutLogger return default logger instance
func DefaultStdoutLogger(c context.Context, logLevel string) logger.Logger {
	return NewEchoGCLogger(c, "default", logLevel, os.Stdout)
}

// NewComponentLogger return default logger instance with custom component
func NewComponentLogger(c context.Context, component string) logger.Logger {
	cfg := logger.Config{}
	err := configurator.NewConfiguration(&cfg, "logger")
	if err != nil {
		panic("Cannot parse config")
	}

	var output io.Writer
	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else {
		output = os.Stdout
	}

	return NewEchoGCLogger(c, component, cfg.LogLevel, output)
}

//nolint:gochecknoinits, reassign
package logger

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/webdevelop-pro/go-common/configurator"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// Printf is implementation of fx.Printer
func (l *Logger) Printf(s string, args ...interface{}) {
	l.Info().Msgf(s, args...)
}

// NewLogger return logger instance
func NewLogger(ctx context.Context, component string, logLevel string, output io.Writer) Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	l := zerolog.
		New(output).
		Level(level).
		Hook(SeverityHook{}).
		Hook(ContextHook{}).
		With().Timestamp()

	if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
		l = l.Caller()
	}

	if ctx != nil {
		l = l.Ctx(ctx)
	}

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
func DefaultStdoutLogger(c context.Context, logLevel string) Logger {
	return NewLogger(c, "default", logLevel, os.Stdout)
}

// NewComponentLogger return default logger instance with custom component
func NewComponentLogger(c context.Context, component string) Logger {
	cfg := Config{}
	err := configurator.NewConfiguration(&cfg, "log")
	if err != nil {
		panic(err)
	}

	var output io.Writer
	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else {
		output = os.Stdout
	}

	return NewLogger(c, component, cfg.LogLevel, output)
}

// FromCtx return default logger instance with custom component
func FromCtx(ctx context.Context, component string) *zerolog.Logger {
	log := zerolog.Ctx(ctx)

	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("component", component)
	})

	return log
}

// NewDefaultLogger return default logger instance
func NewDefaultLogger() Logger {
	cfg := Config{}
	err := configurator.NewConfiguration(&cfg, "log")
	if err != nil {
		panic(err)
	}

	var output io.Writer
	// Beautiful output
	if cfg.LogConsole {
		output = zerolog.NewConsoleWriter()
	} else {
		output = os.Stdout
	}

	return NewLogger(context.Background(), "", cfg.LogLevel, output)
}

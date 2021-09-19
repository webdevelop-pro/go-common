package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

// Logger is wrapper struct around zerolog.Logger that adds some custom functionality
type Logger struct {
	zerolog.Logger
}

// Params ...
type Params struct {
	LogLevel   string    `required:"true" split_words:"true"`
	AppVersion string    `ignored:"true"`
	Component  string    `ignored:"true"`
	output     io.Writer `ignored:"true"`
}

// severityHook structure for add severity field in log
type severityHook struct{}

// NewLogger return logger instance
func NewLogger(params Params) (Logger, error) {
	defaultLogger := GetDefaultLogger(params.output)

	level, err := zerolog.ParseLevel(params.LogLevel)
	if err != nil {
		return defaultLogger, err
	}

	var emptyVarsList []string
	if params.AppVersion == "" {
		emptyVarsList = append(emptyVarsList, "Version")
	}
	if params.Component == "" {
		emptyVarsList = append(emptyVarsList, "Component")
	}

	if len(emptyVarsList) > 0 {
		return defaultLogger, fmt.Errorf("this vars didn't set: %v", emptyVarsList)
	}

	return Logger{
		defaultLogger.Logger.
			Level(level).
			With().
			Str("version", params.AppVersion).
			Str("component", params.Component).
			Logger(),
	}, nil
}

// GetDefaultLogger ...
func GetDefaultLogger(w io.Writer) Logger {
	output := w
	if output == nil {
		output = os.Stdout
	}

	return Logger{
		zerolog.
			New(output).
			Level(zerolog.InfoLevel).
			Hook(severityHook{}).
			With().Caller().
			Timestamp().
			Logger(),
	}
}

// Run convert info to severity
func (h severityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", strings.ToUpper(level.String()))
	} else {
		e.Str("severity", strings.ToUpper(zerolog.ErrorLevel.String()))
		e.Str("message", "Don't use logs with NoLevel")
	}
}

// Printf is implementation of fx.Printer
func (l Logger) Printf(s string, args ...interface{}) {
	l.Info().Msgf(s, args...)
}

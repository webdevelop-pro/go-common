package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/configurator"
)

var (
	cfg           = Config{}
	defaultParams Params
	defaultOutput io.Writer
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := configurator.NewConfiguration(&cfg); err != nil {
		panic(err)
	}

	defaultOutput = os.Stdout

	if cfg.LogConsole {
		defaultOutput = zerolog.NewConsoleWriter()
	}

	defaultParams = Params{
		LogLevel:   cfg.LogLevel,
		AppVersion: "",
		Component:  "",
		output:     defaultOutput,
	}
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

// severityHook is a structure for adding the severity field in log
type severityHook struct{}

// New returns logger instance
func New(params Params) Logger {
	defaultLogger := getDefaultLogger(params.output)

	level, err := zerolog.ParseLevel(params.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
		defaultLogger.Error().Err(err).Str("level_from_params", params.LogLevel).Msg("failed to parse log level")
	}

	ctxLogger := defaultLogger.Level(level).With()

	var emptyVarsList []string
	if params.AppVersion == "" {
		emptyVarsList = append(emptyVarsList, "Version")
	} else {
		ctxLogger = ctxLogger.Str("version", params.AppVersion)
	}

	if params.Component == "" {
		emptyVarsList = append(emptyVarsList, "Component")
	} else {
		ctxLogger = ctxLogger.Str("component", params.Component)
	}

	if len(emptyVarsList) > 0 {
		defaultLogger.Error().
			Err(err).Msg(fmt.Sprintf("this vars didn't set: %v", emptyVarsList))
	}

	return Logger{ctxLogger.Logger()}
}

// GetDefaultParams return default params
func GetDefaultParams() Params {
	return defaultParams
}

func getDefaultLogger(w io.Writer) zerolog.Logger {
	output := w
	if output == nil {
		output = os.Stdout
	}

	return zerolog.
		New(output).
		Level(zerolog.InfoLevel).
		Hook(severityHook{}).
		With().Caller().
		Timestamp().
		Logger()
}

// GetDefaultLogger returns default logger
func GetDefaultLogger(w io.Writer) Logger {
	return Logger{
		getDefaultLogger(w),
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

// NewDefault returns default logger instance
func NewDefault() Logger {
	return New(GetDefaultParams())
}

// NewDefaultComponent returns default logger instance with custom component name
func NewDefaultComponent(component string) Logger {
	params := GetDefaultParams()
	params.Component = component

	return New(params)
}

package logger

import "github.com/rs/zerolog"

const (
	errorTypeKey   = "@type"
	errorTypeValue = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type TypeHook struct {
	skip bool
}

func (h TypeHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if h.skip {
		return
	}

	switch level {
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		e.Str(errorTypeKey, errorTypeValue)
	}
}

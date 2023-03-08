package logger

import "github.com/rs/zerolog"

var defaultSeverityMap = map[zerolog.Level]string{
	zerolog.TraceLevel: "DEBUG",
	zerolog.DebugLevel: "DEBUG",
	zerolog.InfoLevel:  "INFO",
	zerolog.WarnLevel:  "WARNING",
	zerolog.ErrorLevel: "ERROR",
	zerolog.FatalLevel: "CRITICAL",
	zerolog.PanicLevel: "ALERT",
}

type SeverityHook struct{}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if level != zerolog.NoLevel {
		e.Str("severity", defaultSeverityMap[level])
	}
}

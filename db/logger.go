package db

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	comLogger "github.com/webdevelop-pro/go-logger"
)

// Logger is a struct that represent logger for DB
type Logger struct {
	log comLogger.Logger
}

// NewDBLogger is a constructor for Logger
func NewDBLogger(log comLogger.Logger) *Logger {
	return &Logger{log: log}
}

// Log prints a message
func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case tracelog.LogLevelTrace:
		l.log.Trace().Interface("data", data).Msg(msg)
	case tracelog.LogLevelDebug:
		l.log.Debug().Interface("data", data).Msg(msg)
	case tracelog.LogLevelInfo:
		l.log.Info().Interface("data", data).Msg(msg)
	case tracelog.LogLevelWarn:
		l.log.Warn().Interface("data", data).Msg(msg)
	case tracelog.LogLevelError:
		l.log.Error().Interface("data", data).Msg(msg)
	}
}

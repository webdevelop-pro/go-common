package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	comLogger "github.com/webdevelop-pro/go-common/logger"
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
func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case pgx.LogLevelTrace:
		l.log.Trace().Interface("data", data).Msg(msg)
	case pgx.LogLevelDebug:
		l.log.Debug().Interface("data", data).Msg(msg)
	case pgx.LogLevelInfo:
		l.log.Info().Interface("data", data).Msg(msg)
	case pgx.LogLevelWarn:
		l.log.Warn().Interface("data", data).Msg(msg)
	case pgx.LogLevelError:
		l.log.Error().Interface("data", data).Msg(msg)
	}
}

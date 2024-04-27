package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
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
	var err error
	errMsg, ok := data["err"]
	if ok {
		switch val := errMsg.(type) {
		case *pgconn.PgError:
			msg = val.Message
			err = val
		default:
			msg = fmt.Sprintf("%s", val)
			err = errMsg.(error)
		}
	}

	switch level {
	case tracelog.LogLevelNone:
	case tracelog.LogLevelTrace:
		l.log.Trace().Ctx(ctx).Interface("data", data).Msg(msg)
	case tracelog.LogLevelDebug:
		l.log.Debug().Ctx(ctx).Interface("data", data).Msg(msg)
	case tracelog.LogLevelInfo:
		l.log.Info().Ctx(ctx).Interface("data", data).Msg(msg)
	case tracelog.LogLevelWarn:
		l.log.Warn().Ctx(ctx).Interface("data", data).Msg(msg)
	case tracelog.LogLevelError:
		l.log.Error().Ctx(ctx).Interface("data", data).Msg(msg)
	}
}

// LogQuery custom method to log SQL, for future use in Log
func (db *DB) LogQuery(ctx context.Context, query string, args interface{}) {
	// ToDo
	// Replace $1,$2 with values
	q := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(query, "\t", " "),
				"  ", " "),
			"  ", " "),
		"\n", " ")

	db.Log.Trace().Ctx(ctx).Msgf("query: %s, %v", q, args)
}

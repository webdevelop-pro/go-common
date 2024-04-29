package db

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-logger"
)

func GetConfigConn(logger logger.Logger) *pgx.ConnConfig {
	cfg := Config{}

	err := configurator.NewConfiguration(&cfg, pkgName)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot parse config")
	}

	pgConfig, err := pgx.ParseConfig(GetConnString(&cfg))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	pgxLogLevel, err := tracelog.LogLevelFromString(strings.ToLower(cfg.LogLevel))
	if err != nil {
		pgxLogLevel = tracelog.LogLevelError
	}

	pgConfig.Tracer = &tracelog.TraceLog{
		Logger:   NewDBLogger(logger),
		LogLevel: pgxLogLevel,
	}

	return pgConfig
}

// NewConn is constructor for *pgx.Conn
func NewConn() *pgx.Conn {
	logger := logger.NewComponentLogger(pkgName, nil)

	return newConn(GetConfigConn(logger), logger)
}

// NewConnFromConfig is constructor for *pgx.Conn
func NewConnFromConfig(pgConfig *pgx.ConnConfig) *pgx.Conn {
	logger := logger.NewComponentLogger(pkgName, nil)

	return newConn(pgConfig, logger)
}

func newConn(pgConfig *pgx.ConnConfig, logger logger.Logger) *pgx.Conn {
	// ToDo
	// Accept context as parameter
	var (
		pg  *pgx.Conn
		err error
	)

	// we need this to work correctly with GCP
	i := 0
	ticker := time.NewTicker(time.Second)

	for ; ; <-ticker.C {
		i++
		// ToDo, show database type and host
		// logger.Info().Msgf("Connecting to %s, attempt %d", cfg.Type, i)
		logger.Info().Msgf("Connecting to db attempt %d", i)
		pg, err = pgx.ConnectConfig(context.TODO(), pgConfig)
		if err == nil || i > 60 {
			break
		}
		logger.Error().Err(err).Msgf("failed to connect, attempt %d", i)
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	ticker.Stop()
	return pg
}

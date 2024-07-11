package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"
)

// NewConn is constructor for *pgx.Conn
func NewConn(ctx context.Context) *pgx.Conn {
	logger := logger.NewComponentLogger(ctx, pkgName)

	return newConn(ctx, GetConfigConn(logger), logger)
}

// NewConnFromConfig is constructor for *pgx.Conn
func NewConnFromConfig(ctx context.Context, pgConfig *pgx.ConnConfig) *pgx.Conn {
	//nolint:contextcheck
	logger := logger.NewComponentLogger(context.TODO(), pkgName)

	return newConn(ctx, pgConfig, logger)
}

func newConn(ctx context.Context, pgConfig *pgx.ConnConfig, logger logger.Logger) *pgx.Conn {
	// ToDo
	// Accept context as parameter
	var (
		pg  *pgx.Conn
		err error
	)

	i := 0
	ticker := time.NewTicker(time.Second)

	for ; ; <-ticker.C {
		i++
		// ToDo
		// use exponentional backoff connection
		logger.Info().Msgf("Connecting to db attempt %d", i)
		pg, err = pgx.ConnectConfig(ctx, pgConfig)
		if err == nil || i > maxRetries {
			_, err = pg.Exec(ctx, "SET TIME ZONE 'UTC';")
			if err != nil {
				logger.Error().Err(err).Msg("Unable change timezone")
			}
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

func GetConnString(cfg *Config) string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?application_name=%s",
		cfg.Type,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.AppName,
	)
}

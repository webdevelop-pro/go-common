package db

import (
	"context"
	"fmt"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pkg/errors"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"
)

// NewPool is constructor for pgxpool.Pool
func NewPool(ctx context.Context) *pgxpool.Pool {
	logger := logger.NewComponentLogger(ctx, pkgName)
	return newPool(ctx, GetConfigPool(logger), logger)
}

// NewPoolFromConfig is constructor for pgxpool.Pool
func NewPoolFromConfig(ctx context.Context, pgConfig *pgxpool.Config, logger logger.Logger) *pgxpool.Pool {
	return newPool(ctx, pgConfig, logger)
}

func newPool(ctx context.Context, pgConfig *pgxpool.Config, logger logger.Logger) *pgxpool.Pool {
	var pg *pgxpool.Pool
	var err error

	pg, err = backoff.RetryWithData(
		func() (*pgxpool.Pool, error) {
			logger.Debug().Msgf("Connecting to db ")
			return pgxpool.NewWithConfig(ctx, pgConfig)
		},
		backoff.WithContext(
			backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(maxRetries)),
			ctx,
		),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	_, err = pg.Exec(ctx, "SET TIME ZONE 'UTC';")
	if err != nil {
		logger.Error().Err(err).Msg("Unable change timezone")
	}
	return pg
}

func GetConfigPool(logger logger.Logger) *pgxpool.Config {
	cfg := Config{}

	err := configurator.NewConfiguration(&cfg, pkgName)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot parse config")
	}
	pgConfig, err := pgxpool.ParseConfig(GetPoolConnString(&cfg))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	pgxLogLevel, err := tracelog.LogLevelFromString(cfg.LogLevel)
	if err != nil {
		logger.Warn().Err(errors.Wrapf(err, "wrong level: %s", cfg.LogLevel)).Msgf("cannot parse pgxLogLevel")
		pgxLogLevel = tracelog.LogLevelNone
	}

	pgConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   NewDBLogger(logger),
		LogLevel: pgxLogLevel,
	}

	return pgConfig
}

func GetPoolConnString(cfg *Config) string {
	query := GetConnString(cfg)
	return fmt.Sprintf(
		"%s&pool_max_conns=%d&pool_min_conns=%d&pool_max_conn_lifetime=%ds",
		query,
		cfg.MaxConnections,
		cfg.MinConnections,
		cfg.MaxConnLifetime,
	)
}

//nolint:gochecknoglobals
package db

import (
	"context"
	"fmt"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pkg/errors"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-logger"
)

var (
	pkgName    = "db"
	maxRetries = 100
)

func GetConfigPool(logger logger.Logger) *pgxpool.Config {
	cfg := Config{}

	err := configurator.NewConfiguration(&cfg, pkgName)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot parse config")
	}
	pgConfig, err := pgxpool.ParseConfig(GetConnString(&cfg))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	pgConfig.MaxConnLifetime = time.Second * time.Duration(cfg.MaxConnLifetime)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to set conn lifetime")
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

// NewPool is constructor for pgxpool.Pool
func NewPool() *pgxpool.Pool {
	logger := logger.NewComponentLogger(context.TODO(), pkgName)

	return newPool(GetConfigPool(logger), logger)
}

// NewPoolFromConfig is constructor for pgxpool.Pool
func NewPoolFromConfig(pgConfig *pgxpool.Config, logger logger.Logger) *pgxpool.Pool {
	return newPool(pgConfig, logger)
}

func newPool(pgConfig *pgxpool.Config, logger logger.Logger) *pgxpool.Pool {
	var pg *pgxpool.Pool
	var err error
	ctx := context.TODO()

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

	return pg
}

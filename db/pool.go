package db

import (
	"context"
	"fmt"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/webdevelop-pro/go-common/configurator"
	comLogger "github.com/webdevelop-pro/go-logger"
)

var (
	pkgName    = "db"
	logger     = comLogger.NewComponentLogger(pkgName, nil)
	maxRetries = uint64(100)
)

func GetConfigPool(c *configurator.Configurator) *pgxpool.Config {
	cfg := c.New(pkgName, &Config{}, pkgName).(*Config)
	pgConfig, err := pgxpool.ParseConfig(GetConnString(cfg))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	pgConfig.MaxConnLifetime = time.Second * time.Duration(cfg.MaxConnLifetime)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to set conn lifetime")
	}

	pgxLogLevel, err := tracelog.LogLevelFromString(cfg.LogLevel)
	if err != nil {
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
func NewPool(c *configurator.Configurator) *pgxpool.Pool {
	return newPool(GetConfigPool(c))
}

// NewPoolFromConfig is constructor for pgxpool.Pool
func NewPoolFromConfig(pgConfig *pgxpool.Config) *pgxpool.Pool {
	return newPool(pgConfig)
}

func newPool(pgConfig *pgxpool.Config) *pgxpool.Pool {
	var pg *pgxpool.Pool
	var err error
	ctx := context.TODO()

	pg, err = backoff.RetryWithData(
		func() (*pgxpool.Pool, error) {
			logger.Debug().Msgf("Connecting to db ")
			return pgxpool.NewWithConfig(ctx, pgConfig)
		},
		backoff.WithContext(
			backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries),
			ctx,
		),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	return pg
}

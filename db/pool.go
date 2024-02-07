package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/webdevelop-pro/go-common/configurator"
	comLogger "github.com/webdevelop-pro/go-logger"
)

var (
	pkgName = "db"
	logger  = comLogger.NewComponentLogger(pkgName, nil)
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

	// we need this to work correctly with GCP
	i := 0
	ticker := time.NewTicker(time.Second)

	for ; ; <-ticker.C {
		i++
		// ToDo, show database type and host
		// logger.Info().Msgf("Connecting to %s, attempt %d", cfg.Type, i)
		logger.Info().Msgf("Connecting to db attempt %d", i)
		pg, err = pgxpool.NewWithConfig(context.TODO(), pgConfig)
		if err == nil || i > 60 {
			break
		}
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	ticker.Stop()

	return pg
}

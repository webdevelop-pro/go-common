package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/webdevelop-pro/go-common/configurator"
	comLogger "github.com/webdevelop-pro/go-common/logger"
)

var (
	pkgName = "db"
	logger  = comLogger.NewDefaultComponent(pkgName)
)

func ParseConfig(cfg *Config) *pgxpool.Config {
	pgConfig, err := pgxpool.ParseConfig(GetConnString(cfg))

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	pgxLogLevel, err := pgx.LogLevelFromString(cfg.LogLevel)

	if err != nil {
		pgxLogLevel = pgx.LogLevelNone
	}

	pgConfig.ConnConfig.Logger = NewDBLogger(logger)
	pgConfig.ConnConfig.LogLevel = pgxLogLevel

	return pgConfig
}

func GetConfig(c *configurator.Configurator) *Config {
	return c.New(pkgName, &Config{}, pkgName).(*Config)
}

func GetConnString(cfg *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&application_name=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.MaxConnections,
		cfg.HOSTNAME,
	)
}

func GetConnConfig(c *configurator.Configurator) *pgxpool.Config {
	return ParseConfig(GetConfig(c))
}

// NewPool is constructor for pgxpool.Pool
func NewPool(c *configurator.Configurator) *pgxpool.Pool {
	return newPool(GetConnConfig(c))
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
		pg, err = pgxpool.ConnectConfig(context.TODO(), pgConfig)
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

package db

import (
	"context"
	"fmt"
	"net/url"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"
)

// NewConn is constructor for *pgx.Conn
func NewConn(ctx context.Context) (*pgx.Conn, error) {
	log := logger.NewComponentLogger(ctx, pkgName)

	pgConfig, err := GetConfigConn(log)
	if err != nil {
		return nil, err
	}

	return NewConnFromConfig(ctx, pgConfig, log)
}

// NewConnFromConfig is constructor for *pgx.Conn
func NewConnFromConfig(ctx context.Context, pgConfig *pgx.ConnConfig, log logger.Logger) (*pgx.Conn, error) {
	if pgConfig == nil {
		return nil, fmt.Errorf("pgx connection config is nil")
	}

	return newConn(ctx, pgConfig.Copy(), log)
}

func newConn(ctx context.Context, pgConfig *pgx.ConnConfig, log logger.Logger) (*pgx.Conn, error) {
	pg, err := backoff.RetryWithData(
		func() (*pgx.Conn, error) {
			log.Debug().Msg("Connecting to db")
			conn, connectErr := pgx.ConnectConfig(ctx, pgConfig)
			if connectErr != nil {
				log.Error().Err(connectErr).Msg("Failed to connect to db")
				return nil, connectErr
			}

			if setErr := setSessionTimeZone(ctx, conn); setErr != nil {
				_ = conn.Close(ctx)
				return nil, setErr
			}

			return conn, nil
		},
		backoff.WithContext(
			backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(maxRetries)),
			ctx,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	return pg, nil
}

func GetConfigConn(log logger.Logger) (*pgx.ConnConfig, error) {
	cfg, err := configurator.Parse[Config](pkgName)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	pgConfig, err := pgx.ParseConfig(GetConnString(&cfg))
	if err != nil {
		return nil, fmt.Errorf("parse pgx connection config: %w", err)
	}

	configureConnTracing(pgConfig, cfg.LogLevel, log)

	return pgConfig, nil
}

func GetConnString(cfg *Config) string {
	query := url.Values{}
	query.Set("application_name", cfg.AppName)

	return (&url.URL{
		Scheme:   string(cfg.Type),
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:     cfg.Database,
		RawQuery: query.Encode(),
	}).String()
}

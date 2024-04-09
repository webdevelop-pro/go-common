package client_middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/motemen/go-loghttp"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-logger"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func CreateHttpClient(serviceName string, pgPool DB) *http.Client {
	log := logger.NewComponentLogger(serviceName, nil)

	return &http.Client{
		Transport: &loghttp.Transport{
			LogRequest:  logRequest(log, serviceName, pgPool),
			LogResponse: logResponse(log, serviceName, pgPool),
		},
	}
}

func logRequest(log logger.Logger, serviceName string, pgPool DB) func(req *http.Request) {
	return func(req *http.Request) {
		var logID int
		sql := `
			INSERT INTO log_logs(
				service, "type", content_type_id, object_id, path, request_created_at, request_id, request_headers, request_data
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9
			) RETURNING id
		`

		rawBody := []byte("{}")
		if req.Body != nil {
			rawBody, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewReader(rawBody))
		}

		// TODO: Use the same format for incoming logs
		log.Info().Str("path", req.RequestURI).Str("service", serviceName).Msg("Send request to 3rd party")

		err := pgPool.QueryRow(
			req.Context(),
			sql,
			serviceName,
			"outcoming",
			// TODO: я сейчас не знаю как легко связать эти поля, не совсем понимаю для ччего они, нужно обсудить
			1,
			"1",
			req.URL.Path,
			time.Now(),
			req.Context().Value(keys.RequestID).(string),
			req.Header,
			rawBody,
		).Scan(&logID)
		if err != nil {
			log.Warn().Err(err).Msg("can't save log in database")
		}

		*req = *req.WithContext(context.WithValue(req.Context(), keys.RequestLogID, logID))
	}
}

func logResponse(log logger.Logger, serviceName string, pgPool DB) func(req *http.Response) {
	return func(resp *http.Response) {
		logID := resp.Request.Context().Value(keys.RequestLogID).(int)
		sql := `
			UPDATE log_logs
			SET status_code=$2, response_headers=$3, response_data=$4
			WHERE id=$1;
		`

		rawBody := []byte("{}")
		if resp.Body != nil {
			rawBody, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewReader(rawBody))
		}

		// TODO: Use the same format for incoming logs
		log.Trace().Str("path", resp.Request.RequestURI).Str("service", serviceName).Msg("3rd party request finished")

		_, err := pgPool.Exec(
			resp.Request.Context(),
			sql,
			logID,
			resp.StatusCode,
			resp.Header,
			rawBody,
		)
		if err != nil {
			log.Warn().Err(err).Int("log_id", logID).Msg("can't save log in database")
		}
	}
}

package clientmiddleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func CreateHTTPClient(serviceName string, pgPool DB) *http.Client {
	log := logger.NewComponentLogger(context.TODO(), serviceName)

	return &http.Client{
		Transport: &loghttp.Transport{
			LogRequest:  logRequest(log, serviceName, pgPool),
			LogResponse: logResponse(log, serviceName, pgPool), //nolint:bodyclose
		},
	}
}

func getContentID(ctx context.Context, log logger.Logger, model string, pgPool DB) (int, error) {
	contentID := 1

	sql, args, err := sq.Select("id").
		From("django_content_type").Where(sq.And{sq.Eq{"model": model}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return contentID, err
	}

	log.Trace().Msgf("query: %s, %v", sql, args)
	err = pgPool.QueryRow(ctx, sql, args...).Scan(&contentID)
	if err != nil {
		return contentID, err
	}

	return contentID, nil
}

func logRequest(log logger.Logger, serviceName string, pgPool DB) func(req *http.Request) {
	return func(req *http.Request) {
		var (
			logID        int
			reqID, _     = req.Context().Value(keys.RequestID).(string)
			modelType, _ = req.Context().Value(keys.LogObjectType).(string)
			objectID, _  = req.Context().Value(keys.LogObjectID).(int)
			msgID, _     = req.Context().Value(keys.MSGID).(string)
		)

		sql := `
			INSERT INTO log_logs(
				service, "type", content_type_id, object_id, path,
				request_created_at, request_id, request_headers, request_data, msg_id
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			) RETURNING id
		`

		contentID, err := getContentID(req.Context(), log, modelType, pgPool)
		if err != nil {
			log.Warn().Err(err).Msg("can't save log in database")
		}

		rawBody := []byte("{}")
		if req.Body != nil {
			rawBody, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewReader(rawBody))
		}

		// TODO: Use the same format for incoming logs
		log.Info().Str("path", req.RequestURI).Str("service", serviceName).Msg("Send request to 3rd party")

		err = pgPool.QueryRow(
			req.Context(),
			sql,
			serviceName,
			"outcoming",
			// TODO: я сейчас не знаю как легко связать эти поля, не совсем понимаю для ччего они, нужно обсудить
			contentID,
			strconv.Itoa(objectID),
			req.URL.Path,
			time.Now(),
			reqID,
			req.Header,
			rawBody,
			msgID,
		).Scan(&logID)
		if err != nil {
			log.Warn().Err(err).Msg("can't save log in database")
		}

		*req = *req.WithContext(context.WithValue(req.Context(), keys.RequestLogID, logID))
	}
}

func logResponse(log logger.Logger, serviceName string, pgPool DB) func(req *http.Response) {
	return func(resp *http.Response) {
		logID, _ := resp.Request.Context().Value(keys.RequestLogID).(int)

		sql := `
			UPDATE log_logs
			SET status_code=$2, response_headers=$3, response_data=$4
			WHERE id=$1;
		`

		rawBody := []byte("{}")
		if resp.Body != nil {
			rawBody, _ = io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewReader(rawBody))
		}

		if len(string(rawBody)) == 0 {
			rawBody = []byte("{}")
		}

		// TODO: Use the same format for incoming logs
		log.Debug().Str("path", resp.Request.RequestURI).
			Int("logID", logID).Str("service", serviceName).Msg("3rd party request finished")

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

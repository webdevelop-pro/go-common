package db

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors" // Error with stack trace
	"github.com/webdevelop-pro/go-common/configurator"
	comLogger "github.com/webdevelop-pro/go-logger"
)

const ErrNotUpdated = errors.Errorf("UPDATE 0")

// DB is a layer to simplify interact with DB
type DB struct {
	*pgxpool.Pool
	Log comLogger.Logger
}

// New returns new DB instance.
func New(c *configurator.Configurator) *DB {
	return NewDB(NewPool(c), logger)
}

// NewDB returns new DB instance.
func NewDB(pool *pgxpool.Pool, log comLogger.Logger) *DB {
	d := &DB{
		Pool: pool,
		Log:  log,
	}

	return d
}

// Subscribe is
func (db *DB) Subscribe(ctx context.Context, topicName string) (<-chan *[]byte, error) {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(ctx, "listen "+topicName); err != nil {
		return nil, err
	}

	out := make(chan *[]byte)

	go func() {
		defer conn.Release()
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := conn.Conn().WaitForNotification(ctx)
				if err != nil {
					db.Log.Error().Err(err).Msg("Can't receive notification, continuing")

					if conn.Conn().IsClosed() {
						db.Log.Error().Err(err).Msg("Lost connection")
						return
					}

					continue
				}

				payload := []byte(n.Payload)
				out <- &payload
			}
		}
	}()

	return out, nil
}

func (db *DB) LogQuery(ctx context.Context, query string, args interface{}) {
	// ToDo
	// Replace $1,$2 with values
	q := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(query, "\t", " "),
				"  ", " "),
			"  ", " "),
		"\n", " ")

	db.Log.Trace().Ctx(ctx).Msgf("query: %s, %v", q, args)
}

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors" // Error with stack trace
	comLogger "github.com/webdevelop-pro/go-logger"
)

var ErrNotUpdated = errors.Errorf("UPDATE 0")

// DB is a layer to simplify interact with DB
type DB struct {
	*pgxpool.Pool
	Log comLogger.Logger
}

// New returns new DB instance.
func New() *DB {
	logger := comLogger.NewComponentLogger(context.TODO(), pkgName)

	return NewDB(NewPool(), logger)
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

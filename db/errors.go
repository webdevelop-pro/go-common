package db

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	Updated = "UPDATE 0"
	Created = "INSERT 0 1"
)

var (
	ErrNotFound   = fmt.Errorf("db: %w", pgx.ErrNoRows)
	ErrNotUpdated = errors.New("object is not updated")
	ErrSQLQuery   = errors.New("sql query builder error")
	ErrSQLRequest = errors.New("sql request error")
)

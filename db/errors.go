package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	Updated = "UPDATE 0"
	Created = "INSERT 0 1"
)

var (
	ErrNotFound   = errors.Wrapf(pgx.ErrNoRows, "") // so we have stack trace and error message from pgx for errors.Is
	ErrNotUpdated = errors.New("object is not updated")
	ErrSQLQuery   = errors.New("sql query builder error")
	ErrSQLRequest = errors.New("sql request error")
)

package dbtests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/db"
	"github.com/webdevelop-pro/go-common/tests"
)

const (
	sqlRetryInterval = 500
	maxQueryRetries  = 25
)

func getDBPool(t tests.TestContext) *db.DB {
	return t.Ctx.Value("db").(*db.DB)
}

func RawSQL(query string) tests.SomeAction {
	return func(t tests.TestContext) error {
		query = strings.ReplaceAll(query, "\t", " ")
		query = strings.ReplaceAll(query, "\n", " ")
		query = strings.ReplaceAll(query, "  ", " ")

		_, err := getDBPool(t).Pool.Exec(context.Background(), query)
		return err
	}
}

func SQL(query string, expected ...tests.ExpectedResult) tests.SomeAction {
	return func(t tests.TestContext) error {
		var res map[string]interface{}

		query = strings.ReplaceAll(query, "\t", " ")
		query = strings.ReplaceAll(query, "\n", " ")
		query = strings.ReplaceAll(query, "  ", " ")
		rowQuery := "select row_to_json(q)::jsonb from (" + query + ") as q"

		err := getDBPool(t).Pool.QueryRow(context.Background(), rowQuery).Scan(&res)
		// Do XXX retries automatically
		if err != nil {
			try := 0
			ticker := time.NewTicker(sqlRetryInterval * time.Millisecond)
			for range ticker.C {
				err = getDBPool(t).Pool.QueryRow(context.Background(), rowQuery).Scan(&res)
				if err != nil {
					try++
					if try > maxQueryRetries {
						return errors.Wrapf(err, "for sql: %s", query)
					}
				} else {
					break
				}
			}
		}

		for _, exp := range expected {
			for key, value := range exp {
				// ToDo:
				// Find library to have colorful compare for maps
				expValue, ok := res[key]
				if assert.True(t.T, ok, fmt.Sprintf("Expected column %s not exist in result", key)) {
					value = tests.AllowAny(expValue, value)
					assert.Equal(t.T, expValue, value)
				}
			}
		}

		return nil
	}
}

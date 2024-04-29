package tests

import (
	"context"
	"testing"

	"github.com/webdevelop-pro/go-common/db"
)

func TestNewConn(t *testing.T) {
	// ToDo
	// NewConn should accept context
	// ToDo
	// NewConn should return an error
	// ctx := metadata.NewContext(context.Background(), meta)
	conn := db.NewConn()
	var name string
	err := conn.QueryRow(context.Background(), "select 'test'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}

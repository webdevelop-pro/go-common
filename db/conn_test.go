package db

import (
	"context"
	"testing"
)

func TestNewConn(t *testing.T) {
	// ToDo
	// NewConn should return an error
	ctx := context.Background()
	conn := NewConn(ctx)
	var name string
	err := conn.QueryRow(ctx, "select 'test'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}

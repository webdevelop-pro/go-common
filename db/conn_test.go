package db

import (
	"context"
	"testing"
)

func TestNewConn(t *testing.T) {
	ctx := context.Background()
	conn, err := NewConn(ctx)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	var name string
	err = conn.QueryRow(ctx, "select 'test'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}

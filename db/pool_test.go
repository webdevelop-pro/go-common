package db

import (
	"context"
	"testing"
)

func TestNewPool(t *testing.T) {
	ctx := context.Background()
	conn, err := NewPool(ctx)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	var name string
	err = conn.QueryRow(ctx, "select 'test'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}

package db

import (
	"context"
	"testing"
)

func TestNewPool(t *testing.T) {
	ctx := context.Background()
	conn := NewPool(ctx)
	var name string
	err := conn.QueryRow(ctx, "select 'test'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}

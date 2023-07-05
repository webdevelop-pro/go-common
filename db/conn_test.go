package db

import (
	"context"
	"testing"

	"github.com/webdevelop-pro/go-common/configurator"
)

func TestNewConn(t *testing.T) {
	c := configurator.NewConfigurator()

	// ToDo
	// NewConn should accept context
	// ToDo
	// NewConn should return an error
	// ctx := metadata.NewContext(context.Background(), meta)
	conn := NewConn(c)
	var name string
	err := conn.QueryRow(context.Background(), "select 'test'").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}

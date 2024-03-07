package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/context/keys"
)

func TestHTTPCtx(t *testing.T) {
	ctx := context.Background()
	headers := map[string][]string{
		"X-Request-Id":    {"ZXCasdf123"},
		"X-Forwarded-For": {"31.6.1.12"},
	}

	ctx = SetDefaultHTTPCtx(ctx, headers)

	assert.Equal(t, headers["X-Request-Id"][0], keys.GetCtxValue(ctx, keys.RequestID))
	assert.Equal(t, headers["X-Forwarded-For"][0], keys.GetCtxValue(ctx, keys.IPAddress))
}

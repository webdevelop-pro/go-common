package constants

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()

	ctx = SetCtxValue(ctx, RequestID, "RequestID")
	ctx = SetCtxValue(ctx, MSGID, 1)
	ctx = SetCtxValue(ctx, IPAddress, "0.0.0.0")

	assert.Equal(t, "RequestID", GetCtxValue(ctx, RequestID))
	assert.Equal(t, 1, GetCtxValue(ctx, MSGID))
	assert.Equal(t, "0.0.0.0", GetCtxValue(ctx, IPAddress))
	assert.Equal(t, nil, GetCtxValue(ctx, CtxTraceID))
}

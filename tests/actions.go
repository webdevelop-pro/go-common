package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/misc/httputils"
)

type TestContext struct {
	T *testing.T
	//nolint:containedctx
	Ctx context.Context
}

type (
	SomeAction     func(t TestContext) error
	ExpectedResult map[string]interface{}
)

type ExpectedResponse struct {
	Code int
	Body []byte
}

func SendHTTPRequst(req httputils.Request, checks ...ExpectedResponse) SomeAction {
	return func(t TestContext) error {
		res, err := httputils.CreateDefaultRequest(t.Ctx, req)
		assert.NoError(t.T, err)

		result, code, err := httputils.SendRequest(res)
		assert.NoError(t.T, err)

		for _, expected := range checks {
			assert.Equal(t.T, expected.Code, code, "Invalid response code")

			if expected.Body != nil {
				CompareJSONBody(t.T, result, expected.Body)
			}
		}

		return nil
	}
}

func Sleep(d time.Duration) SomeAction {
	return func(_ TestContext) error {
		time.Sleep(d)

		return nil
	}
}

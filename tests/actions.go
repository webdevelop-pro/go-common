package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/webdevelop-pro/go-common/httputils"
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
	Code    int
	Body    []byte
	Headers http.Header
}

func SendHTTPRequst(req httputils.Request, expected ExpectedResponse) SomeAction {
	return func(t TestContext) error {
		newReq, err := httputils.CreateDefaultRequest(t.Ctx, req)
		assert.NoError(t.T, err)

		doRequest := httputils.SendRequest
		if req.HttpClient != nil {
			doRequest = httputils.SendRequestWithClient(req.HttpClient)
		}

		result, headers, code, err := doRequest(newReq)
		assert.NoError(t.T, err)

		assert.Equal(t.T, expected.Code, code, "Invalid response code")
		if expected.Headers != nil {
			assert.EqualValuesf(t.T, expected.Headers, *headers, "headers not equal")
		}

		if expected.Body != nil {
			CompareJSONBody(t.T, result, expected.Body)
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

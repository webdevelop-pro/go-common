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

// DEPRECATED: use SendHTTPRequest instead
func SendHTTPRequst(req httputils.Request, expected ExpectedResponse) SomeAction {
	return SendHTTPRequest(req, expected)
}

func SendHTTPRequest(req httputils.Request, expected ExpectedResponse) SomeAction {
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
			asserts := assert.New(t.T)

			for key := range expected.Headers {
				expectedValue := expected.Headers[key][0]
				actualValue := headers.Get(key)

				if expectedValue == "%any%" {
					continue
				}

				asserts.Equal(expectedValue, actualValue, "Invalid header value for %s", key)
			}
		}

		if expected.Body != nil {
			CompareJSONBody(t.T, result, expected.Body)
		}

		return nil
	}
}

func SendHTTPRequestFiles(req httputils.Request, body map[string]any, files map[string]string, expected ExpectedResponse) SomeAction {
	return func(t TestContext) error {
		newReq, err := httputils.CreateRequestWithFiles(req, body, files)
		assert.NoError(t.T, err)

		doRequest := httputils.SendRequest
		if req.HttpClient != nil {
			doRequest = httputils.SendRequestWithClient(req.HttpClient)
		}

		result, headers, code, err := doRequest(newReq)
		assert.NoError(t.T, err)

		assert.Equal(t.T, expected.Code, code, "Invalid response code")

		if expected.Headers != nil {
			asserts := assert.New(t.T)

			for key := range expected.Headers {
				expectedValue := expected.Headers[key][0]
				actualValue := headers.Get(key)

				if expectedValue == "%any%" {
					continue
				}

				asserts.Equal(expectedValue, actualValue, "Invalid header value for %s", key)
			}
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

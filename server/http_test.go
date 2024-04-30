package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/server/middleware"
	"github.com/webdevelop-pro/go-common/verser"
	logger "github.com/webdevelop-pro/go-logger"
)

func TestHTTPCtx(t *testing.T) {
	ctx := context.Background()
	headers := map[string][]string{
		"X-Request-Id":    {"ZXCasdf123"},
		"X-Forwarded-For": {"31.6.1.12"},
	}

	ctx = middleware.SetDefaultHTTPCtx(ctx, headers)

	assert.Equal(t, headers["X-Request-Id"][0], keys.GetCtxValue(ctx, keys.RequestID))
	assert.Equal(t, headers["X-Forwarded-For"][0], keys.GetCtxValue(ctx, keys.IPAddress))
}

func TestLoggerCtx(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	verser.SetServiVersRepoRevis("test-service", "1.0.0", "gitlab/test/repo", "asdxsgdf")

	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)

	logInfo := logger.ServiceContext{}
	middleware.SetLogger(func(c echo.Context) error {
		ctx := c.Request().Context()

		logInfo, _ = keys.GetCtxValue(ctx, keys.LogInfo).(logger.ServiceContext)

		return nil
	})(echoCtx)

	assert.Equal(t, logInfo.Service, "test-service")
	assert.Equal(t, logInfo.Version, "1.0.0")
	assert.Equal(t, logInfo.SourceReference.Repository, "gitlab/test/repo")
	assert.Equal(t, logInfo.SourceReference.RevisionID, "asdxsgdf")
	assert.Equal(t, logInfo.HTTPRequest.Method, http.MethodPost)
	assert.Equal(t, logInfo.HTTPRequest.URL, "/test")
}

// If Request Id header is empty we should automatically generate it
func TestEmptyRequestID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	echoCtx.Set(echo.HeaderXRequestID, "123123123")
	ctx := middleware.SetDefaultCTX(echoCtx)
	assert.Equal(t, len(keys.GetCtxValue(ctx, keys.RequestID).(string)), 9)
}

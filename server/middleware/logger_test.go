package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/verser"
)

func TestLoggerCtx(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	verser.SetServiVersRepoRevis("test-service", "1.0.0", "gitlab/test/repo", "asdxsgdf")

	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)

	logInfo := logger.ServiceContext{}
	// Set context logger
	SetLogger(func(c echo.Context) error {
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

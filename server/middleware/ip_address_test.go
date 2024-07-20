package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/validator"
)

func TestIPAddress(t *testing.T) {
	// Setup
	e := echo.New()
	// get an instance of a validator
	e.Validator = validator.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	// Set context logger
	SetIPAddress(func(c echo.Context) error {
		return nil
	})(echoCtx)

	assert.Equal(t, keys.GetAsString(echoCtx.Request().Context(), keys.IPAddressStr), "127.0.0.1")
}

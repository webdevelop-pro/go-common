package middleware

import (
	"github.com/labstack/echo/v4"
	logger "github.com/webdevelop-pro/go-logger"
)

type fakeIdentityHeaderMiddleware struct {
	log logger.Logger
}

func NewFakeIdentityHeaderMW() AuthMiddleware {
	l := logger.NewComponentLogger("auth_tool", nil)

	return &fakeIdentityHeaderMiddleware{
		log: l,
	}
}

func (m *fakeIdentityHeaderMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

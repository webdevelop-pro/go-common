package middleware

import (
	"github.com/labstack/echo/v4"
	logger "github.com/webdevelop-pro/go-logger"
)

type FakeIdentityHeaderMiddleware struct {
	log logger.Logger
}

func NewFakeIdentityHeaderMW() *FakeIdentityHeaderMiddleware {
	l := logger.NewComponentLogger(nil, "auth_tool")

	return &FakeIdentityHeaderMiddleware{
		log: l,
	}
}

func (m *FakeIdentityHeaderMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

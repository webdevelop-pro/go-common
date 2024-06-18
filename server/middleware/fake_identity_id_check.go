package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/logger"
)

type FakeIdentityHeaderMiddleware struct {
	log logger.Logger
}

func NewFakeIdentityHeaderMW() *FakeIdentityHeaderMiddleware {
	l := logger.NewComponentLogger(context.TODO(), "auth_tool")

	return &FakeIdentityHeaderMiddleware{
		log: l,
	}
}

func (m *FakeIdentityHeaderMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

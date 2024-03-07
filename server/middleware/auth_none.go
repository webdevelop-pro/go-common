package middleware

import (
	"github.com/labstack/echo/v4"
	logger "github.com/webdevelop-pro/go-logger"
)

type authNoneMiddleware struct {
	log logger.Logger
}

func NewAuthNoneMW() AuthMiddleware {
	l := logger.NewComponentLogger("auth_tool", nil)

	return &authNoneMiddleware{
		log: l,
	}
}

func (m *authNoneMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

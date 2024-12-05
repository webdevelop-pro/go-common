package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/context/keys"

	"github.com/webdevelop-pro/go-common/logger"
)

type AuthIdentityHeaderMiddleware struct {
	log logger.Logger
}

func NewAuthIdentityHeaderMW() *AuthIdentityHeaderMiddleware {
	l := logger.NewComponentLogger(context.TODO(), "auth_tool")

	return &AuthIdentityHeaderMiddleware{
		log: l,
	}
}

func (m *AuthIdentityHeaderMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return CheckIdentityID(next)
}

func CheckIdentityID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		identityID := c.Request().Header.Get("Authorization")

		if identityID == "" {
			return c.JSON(http.StatusUnauthorized, map[string][]string{
				"__error__": {"empty identity_id"},
			},
			)
		}

		ctx := context.WithValue(c.Request().Context(), keys.IdentityID, identityID)
		l := zerolog.Ctx(ctx).With().Str("user_id", identityID).Logger()
		ctx = l.WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

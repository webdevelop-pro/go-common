package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/server/response"
	logger "github.com/webdevelop-pro/go-logger"
)

type authIdentityHeaderMiddleware struct {
	log logger.Logger
}

func NewAuthIdentityHeaderMW() AuthMiddleware {
	l := logger.NewComponentLogger("auth_tool", nil)

	return &authIdentityHeaderMiddleware{
		log: l,
	}
}

func (m *authIdentityHeaderMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		identityID := c.Request().Header.Get("Authorization")

		if identityID == "" {
			return c.JSON(http.StatusUnauthorized, response.Error{
				StatusCode: http.StatusUnauthorized,
				Message: map[string][]string{
					"errors": {"empty identity_id"},
				},
			})
		}

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, "identity_id", identityID)
		l := zerolog.Ctx(ctx).With().Str("user_id", identityID).Logger()
		ctx = l.WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

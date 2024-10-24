package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/httputils"
)

func SetIPAddress(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := httputils.GetIPAddress(c.Request().Header)

		ctx := context.WithValue(c.Request().Context(), keys.IPAddressStr, ip)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

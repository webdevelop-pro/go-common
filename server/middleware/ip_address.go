package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
)

func SetIPAddress(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := keys.GetIPAddress(c.Request().Header)
		c.Set(keys.IPAddressStr, ip)

		//nolint:golint
		ctx := context.WithValue(c.Request().Context(), keys.IPAddressStr, ip)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

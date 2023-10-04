package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// ContextKey is the key used to lookup ip address
	// request from the echo.Context.
	IpAddressContextKey = "ip-address"
)

func SetIPAddress(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := "127.0.0.1"

		if xOFF := c.Request().Header.Get(http.CanonicalHeaderKey("X-Original-Forwarded-For")); xOFF != "" {
			i := strings.Index(xOFF, ", ")
			if i == -1 {
				i = len(xOFF)
			}
			ip = xOFF[:i]
		} else if xFF := c.Request().Header.Get(http.CanonicalHeaderKey("X-Forwarded-For")); xFF != "" {
			i := strings.Index(xFF, ", ")
			if i == -1 {
				i = len(xFF)
			}
			ip = xFF[:i]
		} else if xrIP := c.Request().Header.Get(http.CanonicalHeaderKey("X-Real-IP")); xrIP != "" {
			ip = xrIP
		}

		c.Set(IpAddressContextKey, ip)

		return next(c)
	}
}

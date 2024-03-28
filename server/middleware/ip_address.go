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

func GetIpAddress(headers http.Header) string {
	// ToDo
	// Use echo context.RealIP()
	ip := "127.0.0.1"
	if xOFF := headers.Get("X-Original-Forwarded-For"); xOFF != "" {
		i := strings.Index(xOFF, ", ")
		if i == -1 {
			i = len(xOFF)
		}
		ip = xOFF[:i]
	} else if xFF := headers.Get("X-Forwarded-For"); xFF != "" {
		i := strings.Index(xFF, ", ")
		if i == -1 {
			i = len(xFF)
		}
		ip = xFF[:i]
	} else if xrIP := headers.Get("X-Real-IP"); xrIP != "" {
		ip = xrIP
	}
	return ip
}

func SetIPAddress(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := GetIpAddress(c.Request().Header)
		c.Set(IpAddressContextKey, ip)

		return next(c)
	}
}

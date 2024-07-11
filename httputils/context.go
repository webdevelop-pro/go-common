package httputils

import (
	"net/http"
	"strings"
)

func GetIPAddress(headers http.Header) string {
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

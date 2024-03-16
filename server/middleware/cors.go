package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

var (
	CORS = mw.CORSWithConfig(mw.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			if origin != "" {
				return true, nil
			}

			return false, nil
		},
		AllowCredentials: true,
		AllowMethods:     []string{"GET, POST, PUT, OPTIONS, DELETE, PATCH"},
		AllowHeaders:     []string{"Authorization, X-PINGOTHER, Content-Type, X-Requested-With, X-Request-ID, Vary"},
	})

	CORSHandler = CORS(func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })
)

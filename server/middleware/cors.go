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
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderAuthorization, echo.HeaderContentType,
		},
	})

	CORSHandler = CORS(func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })
)

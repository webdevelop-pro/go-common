package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/response"
)

type ErrorContext struct {
	echo.Context
}

func (c *ErrorContext) ErrorResponse(err error) error {
	log := zerolog.Ctx(c.Request().Context())
	// log := logger.NewComponentLogger("ports-http", c.(logger.Context))
	respErr := response.Error{}
	if errors.As(err, &respErr) {
		if respErr.StatusCode >= http.StatusInternalServerError {
			log.Error().Stack().Err(err).Msgf("system error happen")
		}
		log.Debug().Err(err).Msgf("error response")
		return c.JSON(respErr.StatusCode, respErr.Message)
	}

	// If we have not an response.Error but something else
	log.Warn().Stack().Err(err).Msgf("app return invalid error type")
	resp := map[string]interface{}{"__error__": err.Error()}
	return c.JSON(http.StatusNotImplemented, resp)
}

func (c *ErrorContext) ErrorBadReqestResponse(err error) error {
	log := zerolog.Ctx(c.Request().Context())
	log.Debug().Err(err).Msgf("cannot decode request")

	var resp interface{}
	respErr := response.Error{}
	echoErr := echo.HTTPError{}
	if errors.As(err, &respErr) {
		resp = respErr.Message
	} else if errors.As(err, &echoErr) {
		resp = map[string]interface{}{"__error__": []string{echoErr.Message.(string)}}
	}

	return c.JSON(http.StatusBadRequest, resp)
}

func ErrorHandlers(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &ErrorContext{c}
		return next(cc)
	}
}

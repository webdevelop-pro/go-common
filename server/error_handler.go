package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/response"
)

func ErrorResponse(e echo.Context, err error) error {
	log := zerolog.Ctx(e.Request().Context())
	// log := logger.NewComponentLogger("ports-http", c.(logger.Context))
	respErr := response.Error{}
	if errors.As(err, &respErr) {
		if respErr.StatusCode >= http.StatusInternalServerError {
			log.Error().Stack().Err(err).Msgf("system error happen")
		}
		log.Debug().Err(err).Msgf("error response")
		return e.JSON(respErr.StatusCode, respErr.Message)
	}

	// If we have not an response.Error but something else
	log.Warn().Stack().Err(err).Msgf("app return invalid error type")
	resp := map[string]interface{}{"__error__": err.Error()}
	return e.JSON(http.StatusNotImplemented, resp)
}

func ErrorBadReqestResponse(e echo.Context, err error) error {
	log := zerolog.Ctx(e.Request().Context())
	log.Debug().Err(err).Msgf("cannot decode request")

	var resp interface{}
	respErr := response.Error{}
	echoErr := echo.HTTPError{}
	if errors.As(err, &respErr) {
		resp = respErr.Message
	} else if errors.As(err, &echoErr) {
		resp = map[string]interface{}{"__error__": []string{echoErr.Message.(string)}}
	}

	return e.JSON(http.StatusBadRequest, resp)
}

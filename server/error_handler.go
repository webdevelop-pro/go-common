package server

import (
	"fmt"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/response"
)

func ErrorResponse(e echo.Context, err error) error {
	log := logger.FromCtx(e.Request().Context(), "http/ports")
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
	log := logger.FromCtx(e.Request().Context(), "http/ports")
	log.Debug().Err(err).Msgf("cannot decode request")

	var resp interface{}
	respErr := response.Error{}
	fmt.Println(err.Error())
	if errors.As(err, &respErr) {
		resp = respErr.Message
	} else {
		switch err.(type) {
		case *echo.HTTPError:
			resp = map[string]interface{}{"__error__": []string{err.(*echo.HTTPError).Message.(string)}}
		default:
			resp = map[string]interface{}{"__error__": []string{err.Error()}}
		}
	}

	return e.JSON(http.StatusBadRequest, resp)
}

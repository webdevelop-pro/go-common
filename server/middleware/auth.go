package middleware

import (
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/lib/configurator"
	logger "github.com/webdevelop-pro/lib/logger"
	"github.com/webdevelop-pro/lib/server/errorcode"
	"github.com/webdevelop-pro/lib/server/response"
)

// AuthMiddleware is struct which store instance of auth middleware
type AuthMiddleware struct {
	validateURI string
	log         logger.Logger
}

// Config is a struct to config auth middleware
type Config struct {
	AuthValidateURI string `required:"true" split_words:"true"`
}

// NewAuthMW is a constructor of AuthMiddleware
func NewAuthMW(cfg *Config) *AuthMiddleware {
	return &AuthMiddleware{
		validateURI: cfg.AuthValidateURI,
		log:         logger.NewDefaultComponentLogger("auth_tool"),
	}
}

// NewAuthMiddleware returns a new instance of AuthMiddleware
func NewAuthMiddleware() *AuthMiddleware {
	cfg := &Config{}
	l := logger.NewDefaultComponentLogger("auth_tool")

	if err := configurator.NewConfiguration(cfg); err != nil {
		l.Fatal().Err(err).Msg("failed to get configuration of db")
	}

	return &AuthMiddleware{
		validateURI: cfg.AuthValidateURI,
		log:         l,
	}
}

// Validate is middleware that extracts data from Authorization header and validates it
func (m *AuthMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx   = c.Request().Context()
			l     = zerolog.Ctx(ctx)
			trcID = GetTraceID(ctx)
			// Authorization: Bearer <token>
			token = ExtractTokenFromString(c.Request().Header.Get("Authorization"))
		)

		// make request to auth service
		req, err := http.NewRequest(http.MethodGet, m.validateURI, nil)
		if err != nil {
			l.Error().Err(err).Interface("req", req).Msg("Couldn't form request")
			return c.JSON(http.StatusInternalServerError, response.Error{
				TraceID:     trcID,
				Code:        errorcode.InternalError,
				Description: `couldn't check authenticity`,
			})
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			l.Error().Err(err).Interface("req", req).Msg("Couldn't do request")
			return c.JSON(http.StatusInternalServerError, response.Error{
				TraceID:     trcID,
				Code:        errorcode.InternalError,
				Description: `couldn't check authenticity`,
			})
		}

		// if status code is not 2xx
		if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
			return c.JSON(http.StatusUnauthorized, response.Error{
				TraceID:     trcID,
				Code:        errorcode.BadAuth,
				Description: `not valid token in Authorization header`,
			})
		}

		jwtPayload, err := ParseJWTPayload(token)
		if err != nil {
			l.Error().Err(err).Interface("req", req).Msg("failed to decode token")
			return c.JSON(http.StatusInternalServerError, response.Error{
				TraceID:     trcID,
				Code:        errorcode.InternalError,
				Description: `failed to decode token`,
			})
		}

		if jwtPayload.UserID == "" {
			return c.JSON(http.StatusNotFound, response.Error{
				TraceID:     trcID,
				Code:        errorcode.NoEntity,
				Description: `No such user exists`,
			})
		}

		SetJWTPayload(c, jwtPayload)

		return next(c)
	}
}

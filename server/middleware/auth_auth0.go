package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/webdevelop-pro/go-common/configurator"
	logger "github.com/webdevelop-pro/go-logger"
)

type AuthMiddleware interface {
	Validate(next echo.HandlerFunc) echo.HandlerFunc
}

// AuthMiddleware is struct which store instance of auth middleware
type Auth0Middleware struct {
	validateURI string
	log         logger.Logger
}

// Config is a struct to config auth middleware
type Config struct {
	AuthValidateURI string `required:"true" split_words:"true"`
}

// NewAuthMW is a constructor of AuthMiddleware
func NewAuth0MW(cfg *Config) *Auth0Middleware {
	return &Auth0Middleware{
		validateURI: cfg.AuthValidateURI,
		log:         logger.NewComponentLogger(context.TODO(), "auth_tool"),
	}
}

// NewAuthMiddleware returns a new instance of AuthMiddleware
func NewAuthMiddleware() *Auth0Middleware {
	cfg := &Config{}
	l := logger.NewComponentLogger(context.TODO(), "auth_tool")

	if err := configurator.NewConfiguration(cfg); err != nil {
		l.Fatal().Err(err).Msg("failed to get configuration of db")
	}

	return &Auth0Middleware{
		validateURI: cfg.AuthValidateURI,
		log:         l,
	}
}

// ToDo
// Transfer headers ..
// Validate is middleware that extracts data from Authorization header and validates it
func (m *Auth0Middleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx = c.Request().Context()
			l   = zerolog.Ctx(ctx)
			// Authorization: Bearer <token>
			token = ExtractTokenFromString(c.Request().Header.Get("Authorization"))
		)

		// make request to auth service
		req, err := http.NewRequest(http.MethodGet, m.validateURI, nil)
		if err != nil {
			l.Error().Err(err).Interface("req", req).Msg("Couldn't form request")
			return c.JSON(http.StatusForbidden, map[string][]string{"__error__": {"couldn't check authenticity"}})
		}
		req = req.WithContext(ctx)

		req.Header.Add("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			l.Error().Err(err).Interface("req", req).Msg("Couldn't do request")
			return c.JSON(http.StatusForbidden, map[string][]string{"__error__": {"couldn't check authenticity"}})
		}

		// if status code is not 2xx
		if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
			return c.JSON(
				http.StatusUnauthorized,
				map[string][]string{"__error__": {"not valid token in Authorization header"}},
			)
		}

		jwtPayload, err := ParseJWTPayload(token)
		if err != nil {
			l.Error().Err(err).Interface("req", req).Msg("failed to decode token")
			return c.JSON(
				http.StatusBadRequest,
				map[string][]string{"__error__": {"failed to decode token"}},
			)
		}

		if jwtPayload.UserID == "" {
			return c.JSON(
				http.StatusNotFound,
				map[string][]string{"__error__": {"wrong user"}},
			)
		}

		SetJWTPayload(c, jwtPayload)

		return next(c)
	}
}

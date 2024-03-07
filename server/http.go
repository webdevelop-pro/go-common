package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/server/healthcheck"
	"github.com/webdevelop-pro/go-common/server/middleware"
	"github.com/webdevelop-pro/go-common/server/validator"
	logger "github.com/webdevelop-pro/go-logger"
	"go.uber.org/fx"
)

type HttpServer struct {
	Echo     *echo.Echo
	log      logger.Logger
	config   *Config
	authTool middleware.AuthMiddleware
}

// Route is a http route.
type Route struct {
	Method       string
	Path         string
	Handle       echo.HandlerFunc
	NoCORS       bool
	NoAuth       bool
	IdentityAuth bool
	OptionalAuth bool
	Middlewares  []echo.MiddlewareFunc
}

// AddRoute adds route to the router.
func (h *HttpServer) AddRoute(route Route) {
	handle := route.Handle

	if h.authTool != nil && !route.NoAuth {
		route.Middlewares = append(route.Middlewares, h.authTool.Validate)
	}

	h.Echo.Add(route.Method, route.Path, handle, route.Middlewares...)
}

// SetAuthMiddleware sets auth middleware to the router.
func (h *HttpServer) SetAuthMiddleware(authTool middleware.AuthMiddleware) {
	h.authTool = authTool
}

// NewHttpServer returns new API instance.
func NewHttpServer(e *echo.Echo, l logger.Logger, cfg *Config, authTool middleware.AuthMiddleware) *HttpServer {
	// sets CORS headers if Origin is present
	e.Use(
		echoMW.CORSWithConfig(echoMW.CORSConfig{
			Skipper: func(c echo.Context) bool {
				return false
			},
			AllowOriginFunc: func(origin string) (bool, error) {
				return true, nil
			},
			AllowCredentials: true,
			AllowMethods:     []string{"GET, POST, PUT, OPTIONS, DELETE, PATCH"},
			AllowHeaders:     []string{"Authorization, X-PINGOTHER, Content-Type, X-Requested-With, X-Request-ID, Vary"},
		}),
	)

	// Set context logger
	e.Use(middleware.SetLogger)
	e.Use(middleware.SetRequestTime)
	e.Use(middleware.SetIPAddress)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), keys.RequestID, requestID)))

			return next(c)
		}
	})
	e.Use(middleware.LogRequests)
	// Trace ID middleware generates a unique id for a request.
	e.Use(echoMW.RequestIDWithConfig(echoMW.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, requestID string) {
			c.Set(echo.HeaderXRequestID, requestID)
		},
	}))
	// Add the healthcheck endpoint
	e.GET(`/healthcheck`, healthcheck.Healthcheck)

	// get an instance of a validator
	e.Validator = validator.New()

	// Add prometheus metrics
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Set docs middleware
	// setDocsMiddleware(e)

	// avoid any native logging of echo, because we use custom library for logging
	e.HideBanner = true        // don't log the banner on startup
	e.HidePort = true          // hide log about port server started on
	e.Logger.SetLevel(log.OFF) // disable echo#Logger

	return &HttpServer{
		Echo:     e,
		config:   cfg,
		log:      l,
		authTool: authTool,
	}
}

func New(authTool middleware.AuthMiddleware) *HttpServer {
	cfg := &Config{}
	l := logger.NewComponentLogger("http_server", nil)

	if err := configurator.NewConfiguration(cfg); err != nil {
		l.Fatal().Err(err).Msg("failed to get configuration of server")
	}

	return NewHttpServer(echo.New(), l, cfg, authTool)
}

// StartServer is function that registers start of http server in lifecycle
func StartServer(lc fx.Lifecycle, srv *HttpServer) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				on := fmt.Sprintf("%s:%s", srv.config.Host, srv.config.Port)
				srv.log.Info().Msgf("starting server on %s", on)
				go func() {
					err := srv.Echo.Start(on)
					if err != nil {
						srv.log.Info().Err(err).Msgf("stop server %s", on)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				err := srv.Echo.Shutdown(ctx)
				if err != nil {
					srv.log.Info().Err(err).Msg("couldn't stop server")
				}

				return nil
			},
		},
	)
}

func SetDefaultHTTPCtx(ctx context.Context, headers http.Header) context.Context {
	requestID := headers.Get(echo.HeaderXRequestID)
	IP := middleware.GetIpAddress(headers)

	ctx = keys.SetCtxValue(ctx, keys.RequestID, requestID)
	ctx = keys.SetCtxValue(ctx, keys.IPAddress, IP)
	return ctx
}

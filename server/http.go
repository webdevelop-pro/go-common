package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/server/healthcheck"
	"github.com/webdevelop-pro/go-common/server/middleware"
	"github.com/webdevelop-pro/go-common/server/route"
	"github.com/webdevelop-pro/go-common/server/validator"
	logger "github.com/webdevelop-pro/go-logger"
	"go.uber.org/fx"
)

type HTTPServer struct {
	Echo   *echo.Echo
	log    logger.Logger
	config *Config
}

func InitAllRoutes(srv *HTTPServer, params route.ConfiguratorIn) {
	for _, rg := range params.Configurators {
		srv.InitRoutes(rg)
	}
}

// AddRoute adds route to the router.
func (s *HTTPServer) AddRoute(rte *route.Route) {
	handle := rte.Handle
	rte.Middlewares = append(rte.Middlewares, middleware.SetLogger)
	s.Echo.Add(rte.Method, rte.Path, handle, rte.Middlewares...)
}

func (s *HTTPServer) InitRoutes(rg route.Configurator) {
	for _, rte := range rg.GetRoutes() {
		//nolint:gosec,scopelint
		s.AddRoute(&rte)
	}
}

// NewHTTPServer returns new API instance.
func NewHTTPServer(e *echo.Echo, l logger.Logger, cfg *Config) *HTTPServer {
	// sets CORS headers if Origin is present
	e.Use(
		echoMW.CORSWithConfig(echoMW.CORSConfig{
			Skipper: func(_ echo.Context) bool {
				return false
			},
			AllowOriginFunc: func(_ string) (bool, error) {
				return true, nil
			},
			AllowCredentials: true,
			AllowMethods:     []string{"GET, POST, PUT, OPTIONS, DELETE, PATCH"},
			AllowHeaders:     []string{"Authorization, X-PINGOTHER, Content-Type, X-Requested-With, X-Request-ID, Vary"},
		}),
	)

	// Set context logger
	e.Use(middleware.SetIPAddress)
	e.Use(middleware.DefaultCTXValues)
	e.Use(middleware.SetRequestTime)
	e.Use(middleware.SetLogger)
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

	return &HTTPServer{
		Echo:   e,
		config: cfg,
		log:    l,
	}
}

func New() *HTTPServer {
	cfg := &Config{}
	l := logger.NewComponentLogger(context.TODO(), "http_server")

	if err := configurator.NewConfiguration(cfg); err != nil {
		l.Fatal().Err(err).Msg("failed to get configuration of server")
	}

	return NewHTTPServer(echo.New(), l, cfg)
}

// StartServer is function that registers start of http server in lifecycle
func StartServer(lc fx.Lifecycle, srv *HTTPServer) {
	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
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

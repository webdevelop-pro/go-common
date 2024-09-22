package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/context/keys"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/server/healthcheck"
	"github.com/webdevelop-pro/go-common/server/middleware"
	"github.com/webdevelop-pro/go-common/server/route"
	"github.com/webdevelop-pro/go-common/validator"
	"go.uber.org/fx"
)

const component = "http_server"

type HTTPServer struct {
	Echo   *echo.Echo
	log    logger.Logger
	config *Config
}

func InitAndRun() fx.Option {
	return fx.Module(component,
		// Init http server
		fx.Provide(NewServer),
		fx.Invoke(
			//
			AddDefaultMiddlewares,
			// Registration routes and handlers for http server
			InitHandlerGroups,
			// Run HTTP server
			StartServer,
		),
	)
}

func (s *HTTPServer) InitRoutes(rg route.Configurator) {
	for _, r := range rg.GetRoutes() {
		//nolint:gosec,scopelint
		s.AddRoute(&r)
	}
}

// AddRoute adds route to the router.
func (s *HTTPServer) AddRoute(route *route.Route) {
	s.Echo.Add(route.Method, route.Path, route.Handler, route.Middlewares...)
}

// NewServer returns new API instance.
func NewServer() *HTTPServer {
	var (
		cfg = &Config{}
		l   = logger.NewComponentLogger(context.TODO(), component)
	)

	if err := configurator.NewConfiguration(cfg); err != nil {
		l.Fatal().Err(err).Msg("failed to get configuration of server")
	}

	e := echo.New()
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
	e.Use(middleware.SetLogger)
	// Add Error handlers
	e.Use(middleware.ErrorHandlers)
	// Add the healthcheck endpoint
	e.GET(`/healthcheck`, healthcheck.Healthcheck)

	// get an instance of a validator
	e.Validator = validator.New()

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

func AddPrometheus(srv *HTTPServer) {
	srv.Echo.Use(echoprometheus.NewMiddleware(component))
	srv.Echo.GET("/metrics", echoprometheus.NewHandler())
}

func AddDefaultMiddlewares(srv *HTTPServer) {
	// srv.Echo.Use(echoMW.Recover())
	srv.Echo.Use(echoMW.BodyLimit("20M"))
	srv.Echo.Use(middleware.SetIPAddress)
	srv.Echo.Use(middleware.SetRequestTime)
	srv.Echo.Use(echoMW.BodyDumpWithConfig(echoMW.BodyDumpConfig{
		Skipper: middleware.FileAndHealtchCheckSkipper,
		Handler: middleware.BodyDumpHandler,
	}))
	// Trace ID middleware generates a unique id for a request.
	srv.Echo.Use(echoMW.RequestIDWithConfig(echoMW.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, requestID string) {
			c.Set(echo.HeaderXRequestID, requestID)

			ctx := context.WithValue(c.Request().Context(), keys.RequestID, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
		},
	}))
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

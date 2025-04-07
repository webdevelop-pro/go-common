package server

import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/webdevelop-pro/go-common/context/keys"
	"go.uber.org/fx"

	"github.com/webdevelop-pro/go-common/validator"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"

	"github.com/webdevelop-pro/go-common/server/healthcheck"
	"github.com/webdevelop-pro/go-common/server/middleware"
	"github.com/webdevelop-pro/go-common/server/route"
)

const pkgName = "http_server"

type HTTPServer struct {
	Echo   *echo.Echo
	log    logger.Logger
	config *Config
}

func InitAndRun() fx.Option {
	return fx.Module(pkgName,
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
		l   = logger.NewComponentLogger(context.TODO(), pkgName)
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

	if os.Getenv("HTTP_HEALTHCHECK") != "false" {
		// Add the healthcheck endpoint
		e.GET(`/healthcheck`, healthcheck.Healthcheck)
	}

	// get an instance of a validator
	e.Validator = validator.New()

	// avoid any native logging of echo, because we use custom library for logging
	e.HideBanner = true        // don't log the banner on startup
	e.HidePort = true          // hide log about port server started on
	e.Logger.SetLevel(log.OFF) // disable echo#Logger

	newSrv := &HTTPServer{
		Echo:   e,
		config: cfg,
		log:    l,
	}

	// add HTTPErrorHandler
	newSrv.Echo.HTTPErrorHandler = newSrv.httpErrorHandler

	return newSrv
}

func AddDefaultMiddlewares(srv *HTTPServer) {
	limit := os.Getenv("HTTP_BODY_LIMIT")
	if limit == "" {
		limit = "20M"
	}

	srv.Echo.Use(echoMW.BodyLimit(limit))
	srv.Echo.Use(middleware.SetIPAddress)
	srv.Echo.Use(middleware.SetRequestTime)

	// Trace ID middleware generates a unique id for a request.
	srv.Echo.Use(echoMW.RequestIDWithConfig(echoMW.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, requestID string) {
			c.Set(echo.HeaderXRequestID, requestID)

			ctx := context.WithValue(c.Request().Context(), keys.RequestID, requestID)
			c.SetRequest(c.Request().WithContext(ctx))
		},
	}))

	if os.Getenv("HTTP_PROMETHEUS") != "false" {
		srv.Echo.Use(echoprometheus.NewMiddleware(pkgName))
		srv.Echo.GET("/metrics", echoprometheus.NewHandler())
	}

	if os.Getenv("HTTP_BODY_DUMP") != "false" {
		srv.Echo.Use(echoMW.BodyDumpWithConfig(echoMW.BodyDumpConfig{
			Skipper: middleware.FileAndHealtchCheckSkipper,
			Handler: middleware.BodyDumpHandler,
		}))
	}

	if os.Getenv("HTTP_REQUEST_LOGGER") != "false" {
		srv.Echo.Use(echoMW.RequestLoggerWithConfig(echoMW.RequestLoggerConfig{
			LogURI:       true,
			LogStatus:    true,
			LogMethod:    true,
			LogLatency:   true,
			LogURIPath:   true,
			LogError:     true,
			LogRequestID: true,
			HandleError:  true,

			Skipper: middleware.FileAndHealtchCheckSkipper,
			LogValuesFunc: func(c echo.Context, v echoMW.RequestLoggerValues) error {
				srv.log.Info().
					Str("method", v.Method).
					Str("URI", v.URI).
					Int("status", v.Status).
					Str("request_id", v.RequestID).
					Str("latency", v.Latency.String()).
					Msg("request")

				return nil
			},
		}))
	}

	if os.Getenv("HTTP_REQUEST_RECOVER") != "false" {
		srv.Echo.Use(echoMW.RecoverWithConfig(echoMW.RecoverConfig{
			StackSize: 10 << 10, // 10 KB
			LogLevel:  log.ERROR,
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				srv.log.Error().Err(err).Bytes("stacktrace", stack).Msg("panic recover")

				return err
			},
		}))
	}
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

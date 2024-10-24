package route

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Route is a http route.
type Route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Configurator interface {
	GetRoutes() []Route
}

type ConfiguratorIn struct {
	fx.In

	Configurators []Configurator `group:"route_configurator"`
}

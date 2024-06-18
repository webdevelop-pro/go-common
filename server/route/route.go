package route

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Route is a http route.
type Route struct {
	Method      string
	Path        string
	Handle      echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Configurator interface {
	GetRoutes() []Route
}

type ConfiguratorOut struct {
	fx.Out

	RG Configurator `group:"route_configurator"`
}

type ConfiguratorIn struct {
	// fx.In

	Configurators []Configurator `group:"route_configurator"`
}

func NewConfigurator(rg Configurator) ConfiguratorOut {
	return ConfiguratorOut{
		RG: rg,
	}
}

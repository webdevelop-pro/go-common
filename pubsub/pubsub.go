package pubsub

import (
	"context"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/pubsub/broker"
	"github.com/webdevelop-pro/go-logger"
)

const pkgName = "pubsub"

type PubSubListener struct {
	cfg    *Config
	log    logger.Logger
	routes []PubSubRoute
}

type PubSubRoute struct {
	Topic        string
	Subscription string
	Listener     func(ctx context.Context, msg broker.Message) error
	broker       broker.Broker
}

func New(c *configurator.Configurator, routes []PubSubRoute) PubSubListener {
	log := logger.NewComponentLogger("pubsub", nil)
	cfg := Config{}

	err := configurator.NewConfiguration(&cfg, "pubsub")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot parse pubsub config")
	}

	p := PubSubListener{
		log: log,
		cfg: &cfg,
	}

	p.AddRoutes(routes)

	return p
}

func (p PubSubListener) Start() {
	ctx := context.Background()
	for _, b := range p.routes {
		go b.broker.Listen(ctx, b.Listener)
	}
}

func (p *PubSubListener) AddRoutes(routes []PubSubRoute) error {
	for _, route := range routes {
		broker, err := broker.New(context.Background(), broker.Config{
			ServiceAccountCredentials: p.cfg.ServiceAccountCredentials,
			ProjectID:                 p.cfg.ProjectID,
			Topic:                     route.Topic,
			Subscription:              route.Subscription,
		})
		if err != nil {
			return err
		}

		route.broker = broker
		p.routes = append(p.routes, route)
	}

	return nil
}

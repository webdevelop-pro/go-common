package queue

import (
	"context"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/queue/pclient"
	"github.com/webdevelop-pro/go-logger"
)

type PubSubListener struct {
	cfg    *Config
	log    logger.Logger
	routes []PubSubRoute
}

type PubSubRoute struct {
	Topic        string
	Subscription string
	Listener     func(ctx context.Context, msg pclient.Message) error
	Client       pclient.Client
}

func New(c *configurator.Configurator, routes []PubSubRoute) PubSubListener {
	log := logger.NewComponentLogger("pubsub", nil)
	cfg := Config{}

	err := configurator.NewConfiguration(&cfg, "pubsub")
	if err != nil {
		log.Fatal().Stack().Err(err).Msg(pclient.ErrConfigParse.Error())
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
		br := b
		go br.Client.Listen(ctx, br.Listener)
	}
}

func (p *PubSubListener) AddRoutes(routes []PubSubRoute) error {
	for _, route := range routes {
		client, err := pclient.New(context.Background(), pclient.Config{
			ServiceAccountCredentials: p.cfg.ServiceAccountCredentials,
			ProjectID:                 p.cfg.ProjectID,
			Topic:                     route.Topic,
			Subscription:              route.Subscription,
		})
		if err != nil {
			return err
		}

		route.Client = client
		p.routes = append(p.routes, route)
	}

	return nil
}

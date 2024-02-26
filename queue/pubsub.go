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
	// ToDo
	// enum: event, webhook, raw
	Topic            string
	Name             string
	Subscription     string
	WebhooksListener func(ctx context.Context, msg pclient.Webhook) error
	EventsListener   func(ctx context.Context, msg pclient.Event) error
	MsgsListener     func(ctx context.Context, msg pclient.Message) error
	Client           pclient.Client
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
		if br.Name == "webhooks" {
			go br.Client.ListenWebhooks(ctx, br.WebhooksListener)
			continue
		}
		if br.Name == "events" {
			go br.Client.ListenEvents(ctx, br.EventsListener)
			continue
		}
		if br.Name == "messages" {
			go br.Client.ListenRawMsgs(ctx, br.MsgsListener)
			continue
		}
		log := logger.NewComponentLogger("pubsub", nil)
		log.Fatal().Stack().Msgf("topic name %s incorrect", br.Name)
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

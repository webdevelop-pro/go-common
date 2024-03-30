package queue

import (
	"context"
	"time"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/queue/pclient"
	"github.com/webdevelop-pro/go-logger"
)

type PubSubListener struct {
	cfg    *Config
	log    logger.Logger
	routes []PubSubRoute
	client pclient.Client
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
}

func New(c *configurator.Configurator, routes []PubSubRoute) PubSubListener {
	log := logger.NewComponentLogger("pubsub", nil)

	client, err := pclient.New(context.Background())
	if err != nil {
		log.Fatal().Stack().Err(err).Msg(pclient.ErrConfigParse.Error())
	}

	p := PubSubListener{
		log:    log,
		client: client,
	}

	p.AddRoutes(routes)

	return p
}

func (p PubSubListener) Start() {
	ctx := context.Background()
	for _, b := range p.routes {
		br := b

		if br.Name == "webhooks" {
			go func() {
			CONNECT_WEBHOOKS:
				err := p.client.ListenWebhooks(ctx, br.Subscription, br.Topic, br.WebhooksListener)
				if err != nil {
					time.Sleep(2 * time.Second)
					goto CONNECT_WEBHOOKS
				}
			}()
			continue
		}
		if br.Name == "events" {
			go func() {
			CONNECT_EVENTS:
				err := p.client.ListenEvents(ctx, br.Subscription, br.Topic, br.EventsListener)
				if err != nil {
					time.Sleep(2 * time.Second)
					goto CONNECT_EVENTS
				}
			}()
			continue
		}
		if br.Name == "messages" {
			go func() {
			CONNECT_RAW:
				err := p.client.ListenRawMsgs(ctx, br.Subscription, br.Topic, br.MsgsListener)
				if err != nil {
					time.Sleep(2 * time.Second)
					goto CONNECT_RAW
				}
			}()
			continue
		}
		log := logger.NewComponentLogger("pubsub", nil)
		log.Fatal().Stack().Msgf("topic name %s incorrect", br.Name)
	}
}

func (p *PubSubListener) AddRoutes(routes []PubSubRoute) error {
	p.routes = append(p.routes, routes...)

	return nil
}

package queue

import (
	"context"

	"github.com/webdevelop-pro/go-common/queue/pclient"
	logger "github.com/webdevelop-pro/go-logger"
)

type PubSubListener struct {
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

func New(routes []PubSubRoute) PubSubListener {
	log := logger.NewComponentLogger(context.TODO(), "pubsub")

	client, err := pclient.New(context.Background())
	if err != nil {
		log.Fatal().Stack().Err(err).Msg(pclient.ErrConfigParse.Error())
	}

	p := PubSubListener{
		log:    log,
		client: client,
	}

	err = p.AddRoutes(routes)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed add routes")
	}

	return p
}

func (p PubSubListener) Start() {
	ctx := context.Background()
	for _, b := range p.routes {
		br := b

		switch br.Name {
		case "webhooks":
			go func() {
				for {
					err := p.client.ListenWebhooks(ctx, br.Subscription, br.Topic, br.WebhooksListener)
					p.log.Error().Err(err).Msg("ListenWebhooks error")
				}
			}()
		case "events":
			go func() {
				for {
					err := p.client.ListenEvents(ctx, br.Subscription, br.Topic, br.EventsListener)
					p.log.Error().Err(err).Msg("ListenWebhooks error")
				}
			}()
		case "messages":
			go func() {
				for {
					err := p.client.ListenRawMsgs(ctx, br.Subscription, br.Topic, br.MsgsListener)
					p.log.Error().Err(err).Msg("ListenWebhooks error")
				}
			}()
		default:
			p.log.Fatal().Stack().Msgf("topic name %s incorrect", br.Name)
		}
	}
}

func (p *PubSubListener) AddRoutes(routes []PubSubRoute) error {
	p.routes = append(p.routes, routes...)

	return nil
}

package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/queue/pclient"
)

const (
	// reconnectWindow is how long a listener may stay continuously broken
	// (e.g. its subscription was deleted) while attempting to reconnect
	// before we give up and panic so the supervisor restarts the process.
	reconnectWindow = 60 * time.Second
	// reconnectDelay paces reconnect attempts so a hard failure does not
	// busy-loop the CPU and flood the logs.
	reconnectDelay = 2 * time.Second
	// healthyRunTime is how long a listener must have stayed connected for a
	// subsequent failure to count as a transient drop (window reset) rather
	// than part of a continuous outage.
	healthyRunTime = 15 * time.Second
)

type PubSubListener struct {
	log    logger.Logger
	routes []PubSubRoute
	client *pclient.Client
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
		log.Fatal().Stack().Err(err).Msg(ErrAddRoute.Error())
	}

	return p
}

func (p PubSubListener) Start() {
	ctx := context.Background()
	for _, b := range p.routes {
		br := b

		switch br.Name {
		case "webhooks":
			go p.runListener(br.Name, br.Subscription, func() error {
				return p.client.ListenWebhooks(ctx, br.Subscription, br.Topic, br.WebhooksListener)
			})
		case "events":
			go p.runListener(br.Name, br.Subscription, func() error {
				return p.client.ListenEvents(ctx, br.Subscription, br.Topic, br.EventsListener)
			})
		case "messages":
			go p.runListener(br.Name, br.Subscription, func() error {
				return p.client.ListenRawMsgs(ctx, br.Subscription, br.Topic, br.MsgsListener)
			})
		default:
			p.log.Fatal().Stack().Msgf("%s %s", ErrNotCorrectTopic.Error(), br.Name)
		}
	}
}

// runListener keeps a subscription listener alive across disconnects.
//
// pclient already retries connecting for up to pclient.SubscriptionRetryTimeout
// (so a subscription deleted and recreated within that window reconnects
// transparently). If the listener cannot stay connected for a continuous
// reconnectWindow — e.g. the subscription was deleted and never recreated —
// runListener panics so the process is restarted by its supervisor instead of
// silently busy-looping forever.
func (p PubSubListener) runListener(name, subscription string, listen func() error) {
	var failingSince time.Time

	for {
		start := time.Now()
		err := listen()
		if err == nil {
			// Graceful stop (context cancelled).
			return
		}

		// A listener that stayed up for a while before failing is a transient
		// drop, not a missing subscription: restart the outage window.
		if time.Since(start) >= healthyRunTime {
			failingSince = time.Time{}
		}
		if failingSince.IsZero() {
			failingSince = start
		}

		downFor := time.Since(failingSince)
		if downFor >= reconnectWindow {
			p.log.Error().Err(err).
				Str("subscription", subscription).
				Dur("down_for", downFor).
				Msg(ErrListeningError.Error())
			panic(fmt.Sprintf(
				"pubsub %q listener could not reconnect to subscription %q within %s: %v",
				name, subscription, reconnectWindow, err,
			))
		}

		p.log.Warn().Err(err).
			Str("subscription", subscription).
			Dur("down_for", downFor).
			Msg("pubsub listener disconnected, reconnecting")
		time.Sleep(reconnectDelay)
	}
}

func (p *PubSubListener) AddRoutes(routes []PubSubRoute) error {
	p.routes = append(p.routes, routes...)

	return nil
}

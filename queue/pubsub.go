package queue

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/queue/pclient"
)

const (
	// reconnectWindow is how long a listener may stay continuously broken
	// (e.g. its subscription was deleted) while attempting to reconnect
	// before we give up and panic so the supervisor restarts the process.
	// A permanently-dead subscription stays dead, so this only governs how
	// fast we surface it; PUBSUB_RECONNECT_WINDOW overrides it (e.g. the
	// end-to-end test harness churns subscriptions and runs one long-lived
	// worker across the whole suite, where 60s is too tight).
	reconnectWindow = 60 * time.Second
	// reconnectDelay paces reconnect attempts so a hard failure does not
	// busy-loop the CPU and flood the logs.
	reconnectDelay = 2 * time.Second
	// healthyRunTime is how long a listener must have stayed connected for a
	// subsequent failure to count as a transient drop (window reset) rather
	// than part of a continuous outage.
	healthyRunTime = 15 * time.Second
)

// Deduper persists per-(service,message) processing state so that a Pub/Sub
// message redelivered (at-least-once / ack-deadline expiry) or published more
// than once is processed at-most-once in effect. It is intentionally defined
// with primitive parameters only — the queue package stays DB-agnostic; the
// service injects an implementation backed by its own storage.
//
// Claim returns claimed=true when the caller now owns processing of msgID
// (first sight, or a retry of a previously failed attempt). claimed=false
// means another delivery is in-flight or it was already processed — the caller
// should ack and skip. A non-nil error should be treated as retryable (NACK).
type Deduper interface {
	Claim(ctx context.Context, service, topic, msgID string, attempt int) (claimed bool, err error)
	MarkProcessed(ctx context.Context, service, msgID string) error
	MarkFailed(ctx context.Context, service, msgID, lastErr string) error
}

type PubSubListener struct {
	log     logger.Logger
	routes  []PubSubRoute
	client  *pclient.Client
	deduper Deduper
	service string
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
	return newListener(routes, "", nil)
}

// NewWithDeduper is New plus message-level deduplication. service identifies
// this consumer in the dedup store (one row per service+message); d persists
// the claim/processed/failed state. A nil d behaves exactly like New.
func NewWithDeduper(routes []PubSubRoute, service string, d Deduper) PubSubListener {
	return newListener(routes, service, d)
}

func newListener(routes []PubSubRoute, service string, d Deduper) PubSubListener {
	log := logger.NewComponentLogger(context.TODO(), "pubsub")

	client, err := pclient.New(context.Background())
	if err != nil {
		log.Fatal().Stack().Err(err).Msg(pclient.ErrConfigParse.Error())
	}

	p := PubSubListener{
		log:     log,
		client:  client,
		deduper: d,
		service: service,
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
			cb := p.dedupWebhooks(br.Topic, br.WebhooksListener)
			go p.runListener(br.Name, br.Subscription, func() error {
				return p.client.ListenWebhooks(ctx, br.Subscription, br.Topic, cb)
			})
		case "events":
			cb := p.dedupEvents(br.Topic, br.EventsListener)
			go p.runListener(br.Name, br.Subscription, func() error {
				return p.client.ListenEvents(ctx, br.Subscription, br.Topic, cb)
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

// reconnectWindowDur is the effective reconnect-or-panic window: the
// PUBSUB_RECONNECT_WINDOW env override (a Go duration like "5m") if set and
// valid, otherwise the reconnectWindow default.
func reconnectWindowDur() time.Duration {
	if v := os.Getenv("PUBSUB_RECONNECT_WINDOW"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return reconnectWindow
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
		if downFor >= reconnectWindowDur() {
			p.log.Error().Err(err).
				Str("subscription", subscription).
				Dur("down_for", downFor).
				Msg(ErrListeningError.Error())
			panic(fmt.Sprintf(
				"pubsub %q listener could not reconnect to subscription %q within %s: %v",
				name, subscription, reconnectWindowDur(), err,
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

func attemptOf(a *int) int {
	if a == nil {
		return 0
	}
	return *a
}

// dedupEvents wraps an events callback with claim→handle→finalize. When no
// deduper is configured it returns the original callback unchanged.
func (p PubSubListener) dedupEvents(
	topic string,
	fn func(context.Context, pclient.Event) error,
) func(context.Context, pclient.Event) error {
	if p.deduper == nil || fn == nil {
		return fn
	}
	return func(ctx context.Context, msg pclient.Event) error {
		claimed, err := p.deduper.Claim(ctx, p.service, topic, msg.ID, attemptOf(msg.Attempt))
		if err != nil {
			return err // retryable: NACK
		}
		if !claimed {
			// already processed, or another delivery is in-flight — skip & ack
			return nil
		}
		if herr := fn(ctx, msg); herr != nil {
			if merr := p.deduper.MarkFailed(ctx, p.service, msg.ID, herr.Error()); merr != nil {
				p.log.Error().Err(merr).Str("msg_id", msg.ID).Msg("cannot mark message failed")
			}
			return herr // NACK; redelivery re-claims the failed row
		}
		return p.deduper.MarkProcessed(ctx, p.service, msg.ID)
	}
}

// dedupWebhooks is the webhook-payload counterpart of dedupEvents.
func (p PubSubListener) dedupWebhooks(
	topic string,
	fn func(context.Context, pclient.Webhook) error,
) func(context.Context, pclient.Webhook) error {
	if p.deduper == nil || fn == nil {
		return fn
	}
	return func(ctx context.Context, msg pclient.Webhook) error {
		claimed, err := p.deduper.Claim(ctx, p.service, topic, msg.ID, attemptOf(msg.Attempt))
		if err != nil {
			return err
		}
		if !claimed {
			return nil
		}
		if herr := fn(ctx, msg); herr != nil {
			if merr := p.deduper.MarkFailed(ctx, p.service, msg.ID, herr.Error()); merr != nil {
				p.log.Error().Err(merr).Str("msg_id", msg.ID).Msg("cannot mark message failed")
			}
			return herr
		}
		return p.deduper.MarkProcessed(ctx, p.service, msg.ID)
	}
}

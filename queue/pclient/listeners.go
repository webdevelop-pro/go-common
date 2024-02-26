package pclient

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	gpubsub "cloud.google.com/go/pubsub"
)

func (b *Client) getSubscription(ctx context.Context) (*pubsub.Subscription, error) {
	var err error
	if b.client == nil {
		return nil, ErrNotConnected
	}

	if b.topic == nil {
		return nil, ErrTopicNotSet
	}

	ok, err := b.topic.Exists(ctx)
	if !ok {
		b.log.Fatal().Stack().Err(err).Interface("cfg", b.cfg).Msgf(ErrTopicNotExists.Error())
		return nil, fmt.Errorf("%w: %s", ErrTopicNotExists, b.cfg.Topic)
	}

	if err != nil {
		b.log.Fatal().Stack().Err(err).Interface("cfg", b.cfg).Msgf(ErrTopicConnect.Error())
		return nil, fmt.Errorf("%w: %w", ErrTopicConnect, err)
	}

	name := b.topic.ID()
	sub := b.client.Subscription(b.cfg.Subscription)
	ok, err = sub.Exists(ctx)
	if err != nil {
		b.log.Fatal().Stack().Err(err).Interface("name", name).Msgf(ErrConnectSubscription.Error())
		return nil, fmt.Errorf("%w: %w", ErrConnectSubscription, err)
	}
	if !ok {
		b.log.Fatal().Stack().Err(err).Interface("name", b.cfg.Subscription).Msgf(ErrSubscriptionNotExist.Error())
		return nil, fmt.Errorf("%w: %w", ErrSubscriptionNotExist, err)
	}
	return sub, nil
}

func (b *Client) ListenRawMsgs(ctx context.Context, callback func(ctx context.Context, msg Message) error) error {

	sub, err := b.getSubscription(ctx)
	if err != nil {
		return err
	}
	go b.listenRawGoroutine(ctx, callback, sub)
	return nil
}

func (b *Client) listenRawGoroutine(ctx context.Context, callback func(ctx context.Context, msg Message) error, sub *gpubsub.Subscription) error {
	// Start consuming messages from the subscription
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		// Unmarshal the message data into a struct
		m := Message{}
		m.Data = msg.Data
		m.Attributes = msg.Attributes
		m.ID = msg.ID

		b.log.Debug().Str("msg", string(m.Data)).Msgf("received message")
		err := callback(ctx, m)
		if err != nil {
			b.log.Error().Err(err).Msgf(ErrReceiveCallback.Error())
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		b.log.Fatal().Stack().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	b.log.Trace().Msgf("connected to subscription %s listen messages", b.cfg.Subscription)
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen messages for %s", sub.ID())
	return nil
}

func (b *Client) ListenWebhooks(ctx context.Context, callback func(ctx context.Context, msg Webhook) error) error {

	sub, err := b.getSubscription(ctx)
	if err != nil {
		return err
	}
	go b.listenWebhookGoroutine(ctx, callback, sub)
	return nil
}

func (b *Client) listenWebhookGoroutine(ctx context.Context, callback func(ctx context.Context, msg Webhook) error, sub *gpubsub.Subscription) error {
	// Start consuming messages from the subscription
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		webhook := Webhook{}
		if err := json.Unmarshal(msg.Data, &webhook); err != nil {
			b.log.Error().Err(fmt.Errorf("cannot unmarshal: %s", msg.Data)).Msgf(ErrUnmarshalPubSub.Error())
			msg.Nack()
			return
		}
		webhook.ID = msg.ID

		b.log.Debug().Interface("msg", webhook).Msgf("received webhook")
		err := callback(ctx, webhook)
		if err != nil {
			b.log.Error().Err(err).Msgf(ErrReceiveCallback.Error())
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		b.log.Fatal().Stack().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	b.log.Trace().Msgf("connected to subscription %s listen for webhooks", b.cfg.Subscription)
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen for webhooks for %s", sub.ID())
	return nil
}

func (b *Client) ListenEvents(ctx context.Context, callback func(ctx context.Context, msg Event) error) error {

	sub, err := b.getSubscription(ctx)
	if err != nil {
		return err
	}
	go b.listenEventGoroutine(ctx, callback, sub)
	return nil
}

func (b *Client) listenEventGoroutine(ctx context.Context, callback func(ctx context.Context, msg Event) error, sub *gpubsub.Subscription) error {
	// Start consuming messages from the subscription
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		event := Event{}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			b.log.Error().Err(fmt.Errorf("cannot unmarshal: %s", msg.Data)).Msgf(ErrUnmarshalPubSub.Error())
			msg.Nack()
			return
		}
		event.ID = msg.ID

		b.log.Debug().Interface("msg", event).Msgf("received event")
		err := callback(ctx, event)
		if err != nil {
			b.log.Error().Err(err).Msgf(ErrReceiveCallback.Error())
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		b.log.Fatal().Stack().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	b.log.Trace().Msgf("connected to subscription %s listen for events", b.cfg.Subscription)
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen for events for %s", sub.ID())
	return nil
}

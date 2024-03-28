package pclient

import (
	"context"
	"encoding/json"
	"fmt"

	gpubsub "cloud.google.com/go/pubsub"
	"github.com/webdevelop-pro/go-common/context/keys"
)

func (b *Client) getSubscription(ctx context.Context, subscription, topic string) (*gpubsub.Subscription, error) {
	var err error
	if b.client == nil {
		return nil, ErrNotConnected
	}

	bTopic := b.client.Topic(topic)

	if bTopic == nil {
		return nil, ErrTopicNotSet
	}

	ok, err := bTopic.Exists(ctx)
	if !ok {
		b.log.Fatal().Stack().Err(err).Str("topic", topic).Msgf(ErrTopicNotExists.Error())
		return nil, fmt.Errorf("%w: %s", ErrTopicNotExists, bTopic.ID())
	}

	if err != nil {
		b.log.Fatal().Stack().Err(err).Str("topic", topic).Msgf(ErrTopicConnect.Error())
		return nil, fmt.Errorf("%w: %w", ErrTopicConnect, err)
	}

	sub := b.client.Subscription(subscription)
	ok, err = sub.Exists(ctx)
	if err != nil {
		b.log.Fatal().Stack().Err(err).Interface("name", subscription).Msgf(ErrConnectSubscription.Error())
		return nil, fmt.Errorf("%w: %w", ErrConnectSubscription, err)
	}
	if !ok {
		b.log.Fatal().Stack().Err(err).Interface("name", subscription).Msgf(ErrSubscriptionNotExist.Error())
		return nil, fmt.Errorf("%w: %w", ErrSubscriptionNotExist, err)
	}
	return sub, nil
}

func (b *Client) ListenRawMsgs(ctx context.Context, subscription, topic string, callback func(ctx context.Context, msg Message) error) error {
	sub, err := b.getSubscription(ctx, subscription, topic)
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

		ctx = keys.SetCtxValue(ctx, keys.MSGID, msg.ID)
		b.log.Trace().Str("msg", string(m.Data)).Msgf("received message")
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
	b.log.Trace().Msgf("connected to subscription %s listen messages", sub.ID())
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen messages for %s", sub.ID())
	return nil
}

func (b *Client) ListenWebhooks(ctx context.Context, subscription, topic string, callback func(ctx context.Context, msg Webhook) error) error {
	sub, err := b.getSubscription(ctx, subscription, topic)
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
			b.log.Error().Err(err).Interface("data", string(msg.Data)).Msgf(ErrUnmarshalPubSub.Error())
			msg.Nack()
			return
		}
		webhook.ID = msg.ID

		ctx = SetDefaultWebhookCtx(ctx, webhook)
		b.log.Trace().Interface("msg", webhook).Msgf("received webhook")
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
	b.log.Trace().Msgf("connected to subscription %s listen for webhooks", sub.ID())
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen for webhooks for %s", sub.ID())
	return nil
}

func (b *Client) ListenEvents(ctx context.Context, subscription, topic string, callback func(ctx context.Context, msg Event) error) error {
	sub, err := b.getSubscription(ctx, subscription, topic)
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
			b.log.Error().Err(err).Interface("data", string(msg.Data)).Msgf(ErrUnmarshalPubSub.Error())
			msg.Nack()
			return
		}
		event.ID = msg.ID

		ctx = SetDefaultEventCtx(ctx, event)

		b.log.Trace().Interface("msg", event).Msgf("received event")
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
	b.log.Trace().Msgf("connected to subscription %s listen for events", sub.ID())
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen for events for %s", sub.ID())
	return nil
}

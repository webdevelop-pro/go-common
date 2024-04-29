package pclient

import (
	"context"
	"encoding/json"
	"fmt"

	gpubsub "cloud.google.com/go/pubsub"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/webdevelop-pro/go-common/context/keys"
)

const (
	maxRetries         = 100
	maxDeliveryAttempt = 10
)

func verifyDeliveryAttempt(msg *gpubsub.Message) {
	// ToDo
	// For some reason right now message does not goes in dead letter queue
	// Fix dead letter queue settings in GCP
	// For now we just ask message to stop working with it
	if msg.DeliveryAttempt != nil && *msg.DeliveryAttempt > maxDeliveryAttempt {
		msg.Ack()
		return
	}
}

func (b *Client) getSubscriptionRetry(ctx context.Context, subscription, topic string) (*gpubsub.Subscription, error) {
	sub, err := backoff.RetryWithData(
		func() (*gpubsub.Subscription, error) {
			b.log.Info().Msgf("Connecting to subscription %s/%s", topic, subscription)
			return b.getSubscription(ctx, subscription, topic)
		},
		backoff.WithContext(
			backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries),
			ctx,
		),
	)
	if err != nil {
		b.log.Error().Stack().Err(err).Msgf(ErrNotConnected.Error())
		return nil, err
	}
	return sub, nil
}

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
		b.log.Error().Err(err).Str("topic", topic).Msgf(ErrTopicNotExists.Error())
		return nil, fmt.Errorf("%w: %s", ErrTopicNotExists, bTopic.ID())
	}

	if err != nil {
		b.log.Error().Err(err).Str("topic", topic).Msgf(ErrTopicConnect.Error())
		return nil, fmt.Errorf("%w: %w", ErrTopicConnect, err)
	}

	sub := b.client.Subscription(subscription)
	ok, err = sub.Exists(ctx)
	if err != nil {
		b.log.Error().Err(err).Str("subscription", subscription).Msgf(ErrConnectSubscription.Error())
		return nil, fmt.Errorf("%w: %w", ErrConnectSubscription, err)
	}
	if !ok {
		b.log.Error().Err(err).Str("subscription", subscription).Msgf(ErrSubscriptionNotExist.Error())
		return nil, fmt.Errorf("%w: %w", ErrSubscriptionNotExist, err)
	}
	return sub, nil
}

func (b *Client) ListenRawMsgs(
	ctx context.Context,
	subscription, topic string,
	callback func(ctx context.Context, msg Message) error,
) error {
	sub, err := b.getSubscriptionRetry(ctx, subscription, topic)
	if err != nil {
		return err
	}
	return b.listenRawGoroutine(ctx, callback, sub)
}

func (b *Client) listenRawGoroutine(
	ctx context.Context,
	callback func(ctx context.Context, msg Message) error,
	sub *gpubsub.Subscription,
) error {
	// Start consuming messages from the subscription
	b.log.Trace().Msgf("connected to subscription %s listen messages", sub.ID())
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		verifyDeliveryAttempt(msg)

		// Unmarshal the message data into a struct
		m := Message{}
		m.Data = msg.Data
		m.Attributes = msg.Attributes
		m.ID = msg.ID
		m.Attempt = msg.DeliveryAttempt

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
		b.log.Error().Stack().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	return err
}

func (b *Client) ListenWebhooks(
	ctx context.Context, subscription,
	topic string,
	callback func(ctx context.Context, msg Webhook) error,
) error {
	sub, err := b.getSubscriptionRetry(ctx, subscription, topic)
	if err != nil {
		return err
	}

	return b.listenWebhookGoroutine(ctx, callback, sub)
}

func (b *Client) listenWebhookGoroutine(
	ctx context.Context,
	callback func(ctx context.Context, msg Webhook) error,
	sub *gpubsub.Subscription,
) error {
	// Start consuming messages from the subscription
	b.log.Trace().Msgf("connected to subscription %s listen for webhooks", sub.ID())
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		verifyDeliveryAttempt(msg)

		webhook := Webhook{}
		if err := json.Unmarshal(msg.Data, &webhook); err != nil {
			b.log.Error().Err(err).Interface("data", string(msg.Data)).Msgf(ErrUnmarshalPubSub.Error())
			msg.Nack()
			return
		}
		webhook.ID = msg.ID
		webhook.Attempt = msg.DeliveryAttempt

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
		b.log.Error().Stack().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	return err
}

func (b *Client) ListenEvents(
	ctx context.Context,
	subscription,
	topic string,
	callback func(ctx context.Context, msg Event) error,
) error {
	sub, err := b.getSubscriptionRetry(ctx, subscription, topic)
	if err != nil {
		return err
	}
	return b.listenEventGoroutine(ctx, callback, sub)
}

func (b *Client) listenEventGoroutine(
	ctx context.Context,
	callback func(ctx context.Context, msg Event) error,
	sub *gpubsub.Subscription,
) error {
	// Start consuming messages from the subscription
	b.log.Trace().Msgf("connected to subscription %s listen for events", sub.ID())
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		verifyDeliveryAttempt(msg)

		event := Event{}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			b.log.Error().Err(err).Interface("data", string(msg.Data)).Msgf(ErrUnmarshalPubSub.Error())
			msg.Nack()
			return
		}
		event.ID = msg.ID
		event.Attempt = msg.DeliveryAttempt

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
		b.log.Error().Stack().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	return err
}

package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/webdevelop-pro/go-logger"
	"google.golang.org/api/option"

	gpubsub "cloud.google.com/go/pubsub"
)

const pkgName = "pubsub"

type Broker struct {
	client *gpubsub.Client // google cloud pubsub client
	topic  *gpubsub.Topic  // google cloud pubsub topic
	log    logger.Logger   // Broker logger
	cfg    *Config         // broker config
}

func New(ctx context.Context, cfg Config) (Broker, error) {
	var err error

	b := Broker{
		log: logger.NewComponentLogger(pkgName, nil),
		cfg: &cfg,
	}

	b.log.Trace().Msgf("Connecting to %s", b.cfg.Topic)
	b.client, err = gpubsub.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.ServiceAccountCredentials))
	if err != nil {
		b.log.Error().Err(err).Interface("cfg", b.cfg).Msgf(ErrConnection.Error())
		return b, fmt.Errorf("%w: %w", ErrConnection, err)
	}

	b.topic = b.client.Topic(b.cfg.Topic)
	return b, nil
}

func (b *Broker) CreateTopic(ctx context.Context, name string) (*gpubsub.Topic, error) {
	b.log.Trace().Msgf("creating topic %s", name)
	if b.client == nil {
		return nil, ErrNotConnected
	}
	topic, err := b.client.CreateTopic(ctx, name)
	if err != nil {
		b.log.Error().Err(err).Interface("name", name).Msgf(ErrTopicCreate.Error())
		return nil, fmt.Errorf("%w: %w", ErrTopicCreate, err)
	}
	return topic, nil
}

func (b *Broker) DeleteTopic(ctx context.Context, name string) error {
	b.log.Trace().Msgf("deleting topic %s", name)
	if b.client == nil {
		return ErrNotConnected
	}
	topic := b.client.Topic(name)
	return topic.Delete(ctx)
}

func (b *Broker) CreateSubscription(ctx context.Context, name string) (*gpubsub.Subscription, error) {
	b.log.Trace().Msgf("creating subscription %s", name)
	if b.client == nil {
		return nil, ErrNotConnected
	}
	if b.topic == nil {
		return nil, ErrTopicNotSet
	}
	// FixME
	// Add RetryPolicy
	sub, err := b.client.CreateSubscription(ctx, name, gpubsub.SubscriptionConfig{
		Topic: b.topic,
		RetryPolicy: &gpubsub.RetryPolicy{
			MinimumBackoff: time.Duration(2),
			MaximumBackoff: time.Duration(120),
		},
	})
	if err != nil {
		b.log.Error().Err(err).Interface("name", name).Msgf(ErrCreateSubscription.Error())
		return nil, fmt.Errorf("%w: %w", ErrCreateSubscription, err)
	}
	return sub, nil
}

func (b *Broker) SetTopic(topic string) error {
	if b.client == nil {
		return ErrNotConnected
	}
	b.topic = b.client.Topic(topic)
	return nil
}

func (b *Broker) Close() {
	if b.client != nil {
		// timeout for Listen to process all messages
		time.Sleep(time.Second)
		if err := b.client.Close(); err != nil {
			b.log.Error().Err(err).Msgf(ErrCloseConnection.Error())
		}
	}
}

type T interface {
}

func (b *Broker) Listen(ctx context.Context, callback func(ctx context.Context, msg Message) error) error {
	var err error

	if b.client == nil {
		return ErrNotConnected
	}

	if b.topic == nil {
		return ErrTopicNotSet
	}

	ok, err := b.topic.Exists(ctx)

	if !ok {
		b.log.Fatal().Interface("cfg", b.cfg).Msgf(ErrTopicNotExists.Error())
		return fmt.Errorf("%w: %s", ErrTopicNotExists, b.cfg.Topic)
	}

	if err != nil {
		b.log.Fatal().Err(err).Interface("cfg", b.cfg).Msgf(ErrTopicConnect.Error())
		return fmt.Errorf("%w: %w", ErrTopicConnect, err)
	}

	name := b.topic.ID()
	sub := b.client.Subscription(b.cfg.Subscription)
	ok, err = sub.Exists(ctx)
	if err != nil {
		b.log.Fatal().Err(err).Interface("name", name).Msgf(ErrConnectSubscription.Error())
		return fmt.Errorf("%w: %w", ErrConnectSubscription, err)
	}
	if !ok {
		b.log.Fatal().Err(err).Interface("name", name).Msgf(ErrSubscriptionNotExist.Error())
		return fmt.Errorf("%w: %w", ErrSubscriptionNotExist, err)
	}
	b.log.Trace().Msgf("connected to subscription %s listen for new messages", name)
	go b.listenGoroutine(ctx, callback, sub)
	return nil
}

func (b *Broker) listenGoroutine(ctx context.Context, callback func(ctx context.Context, msg Message) error, sub *gpubsub.Subscription) error {
	// Start consuming messages from the subscription
	err := sub.Receive(ctx, func(ctx context.Context, msg *gpubsub.Message) {
		// Unmarshal the message data into a struct
		m := Message{}
		m.Data = msg.Data
		m.Attributes = msg.Attributes
		m.ID = msg.ID

		b.log.Trace().Str("msg", string(m.Data)).Msgf("received new message")
		err := callback(ctx, m)
		if err != nil {
			b.log.Error().Err(err).Msgf(ErrReceiveCallback.Error())
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		b.log.Fatal().Err(err).Msgf(ErrReceiveSubscription.Error())
	}
	<-ctx.Done()
	b.log.Trace().Msgf("stop listen for new messages for %s", sub.ID())
	return nil
}

func (b *Broker) Publish(ctx context.Context, msg *Message) (string, error) {

	if b.topic == nil {
		return "", ErrTopicNotSet
	}

	gMsg := gpubsub.Message{Data: msg.Data, Attributes: msg.Attributes}
	msgId, err := b.topic.Publish(ctx, &gMsg).Get(ctx)

	if err != nil {
		b.log.Error().Err(err).Interface("msg", msg).Msgf(ErrPublish.Error())
		return "", fmt.Errorf("%w: %w", ErrPublish, err)
	}

	b.log.Trace().Msgf("sent message %s to %s", msgId, b.topic.ID())
	return msgId, nil
}

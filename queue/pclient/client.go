package pclient

import (
	"context"
	"fmt"
	"time"

	"github.com/webdevelop-pro/go-logger"
	"google.golang.org/api/option"

	gpubsub "cloud.google.com/go/pubsub"
)

const pkgName = "pubsub"

type Client struct {
	client *gpubsub.Client // google cloud pubsub client
	topic  *gpubsub.Topic  // google cloud pubsub topic
	log    logger.Logger   // client logger
	cfg    *Config         // client config
}

func New(ctx context.Context, cfg Config) (Client, error) {
	var err error

	b := Client{
		log: logger.NewComponentLogger(pkgName, nil),
		cfg: &cfg,
	}

	b.log.Trace().Msgf("New pubsub %s", cfg.Topic)
	b.client, err = gpubsub.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.ServiceAccountCredentials))
	if err != nil {
		b.log.Error().Err(err).Interface("cfg", b.cfg).Msgf(ErrConnection.Error())
		return b, fmt.Errorf("%w: %w", ErrConnection, err)
	}

	b.topic = b.client.Topic(cfg.Topic)
	return b, nil
}

func (b *Client) CreateTopic(ctx context.Context, name string) (*gpubsub.Topic, error) {
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

func (b *Client) DeleteTopic(ctx context.Context, name string) error {
	b.log.Trace().Msgf("deleting topic %s", name)
	if b.client == nil {
		return ErrNotConnected
	}
	topic := b.client.Topic(name)
	return topic.Delete(ctx)
}

func (b *Client) DeleteSubscription(ctx context.Context, name string) error {
	b.log.Trace().Msgf("deleting subscription %s", name)
	if b.client == nil {
		return ErrNotConnected
	}
	subscription := b.client.Subscription(name)
	return subscription.Delete(ctx)
}

func (b *Client) CreateSubscription(ctx context.Context, name string) (*gpubsub.Subscription, error) {
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

func (b *Client) SetTopic(ctx context.Context, topic string) error {
	b.topic = b.client.Topic(topic)
	/*
		ok, err := b.topic.Exists(ctx)
		if !ok {
			b.log.Error().Err(err).Interface("topic", topic).Msgf(ErrTopicNotExists.Error())
			return fmt.Errorf("%w: %s", err, topic)
		}
	*/
	return nil
}

func (b *Client) TopicExist(ctx context.Context, topic string) (bool, error) {
	if b.client == nil {
		return false, ErrNotConnected
	}
	top := b.client.Topic(topic)
	ok, err := top.Exists(ctx)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Client) SubscriptionExist(ctx context.Context, sub string) (bool, error) {
	if b.client == nil {
		return false, ErrNotConnected
	}
	subscription := b.client.Subscription(sub)
	ok, err := subscription.Exists(ctx)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Client) Close() {
	if b.client != nil {
		// timeout for Listen to process all messages
		time.Sleep(time.Second)
		if err := b.client.Close(); err != nil {
			b.log.Error().Err(err).Msgf(ErrCloseConnection.Error())
		}
	}
}

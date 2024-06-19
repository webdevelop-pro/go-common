package pclient

import (
	"context"
	"fmt"
	"time"

	gpubsub "cloud.google.com/go/pubsub"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"
	"google.golang.org/api/option"
)

const (
	pkgName         = "pubsub"
	minPubsubBackof = 5
	maxPubsubBackof = 10
)

type Client struct {
	client *gpubsub.Client // google cloud pubsub client
	log    logger.Logger   // client logger
	cfg    *Config         // client config
}

func New(ctx context.Context) (Client, error) {
	cfg := Config{}
	log := logger.NewComponentLogger(ctx, pkgName)

	err := configurator.NewConfiguration(&cfg, pkgName)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg(ErrConfigParse.Error())
		return Client{}, err
	}

	bclient, err := gpubsub.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.ServiceAccountCredentials))
	if err != nil {
		log.Fatal().Stack().Err(err).Msg(err.Error())
		return Client{}, err
	}

	b := Client{
		log:    log,
		cfg:    &cfg,
		client: bclient,
	}

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

func (b *Client) CreateSubscription(ctx context.Context, name, topic string) (*gpubsub.Subscription, error) {
	b.log.Trace().Msgf("creating subscription %s", name)
	if b.client == nil {
		return nil, ErrNotConnected
	}

	pTopic := b.client.Topic(topic)

	// FixME
	// Add RetryPolicy
	sub, err := b.client.CreateSubscription(ctx, name, gpubsub.SubscriptionConfig{
		Topic:                     pTopic,
		EnableExactlyOnceDelivery: true,
		RetryPolicy: &gpubsub.RetryPolicy{
			MinimumBackoff: time.Minute * minPubsubBackof,
			MaximumBackoff: time.Minute * maxPubsubBackof,
		},
	})
	if err != nil {
		b.log.Error().Err(err).Interface("name", name).Msgf(ErrCreateSubscription.Error())
		return nil, fmt.Errorf("%w: %w", ErrCreateSubscription, err)
	}
	return sub, nil
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

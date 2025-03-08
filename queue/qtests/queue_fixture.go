package qtests

import (
	"context"
	"fmt"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/queue/pclient"
)

type contextKey string

const queueKey contextKey = "queue"

type Fixture struct {
	topic        string
	subscription string
	filePath     string
}

func NewFixture(topic, subscription, filePath string) Fixture {
	return Fixture{
		topic:        topic,
		subscription: subscription,
		filePath:     filePath,
	}
}

type FixturesManager struct {
	queue    *pclient.Client
	cfg      pclient.Config
	fixtures []Fixture
}

func NewFixturesManager(ctx context.Context, fixtures ...Fixture) FixturesManager {
	configurator.LoadDotEnv()
	cfg := pclient.Config{}

	pclient, err := pclient.New(ctx)
	if err != nil {
		panic(err)
	}
	return FixturesManager{
		queue:    pclient,
		cfg:      cfg,
		fixtures: fixtures,
	}
}

func (f FixturesManager) CleanAndApply() error {
	for _, fixture := range f.fixtures {
		err := f.Clean(fixture.topic, fixture.subscription)
		if err != nil {
			return err
		}
	}
	// ToDo
	// Push data to the subscriptions
	// return PubSubF.LoadFixtures(fixture.filePath)
	return nil
}

func (f FixturesManager) SetCTX(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "testm", "test test")
	return context.WithValue(ctx, queueKey, f.queue)
}

func (f FixturesManager) Clean(topic string, subscription string) error {
	ctx := context.Background()
	ok, err := f.queue.SubscriptionExist(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed check subscription: %w", err)
	}
	if ok {
		err := f.queue.DeleteSubscription(ctx, subscription)
		if err != nil {
			return fmt.Errorf("failed delete subscription: %w", err)
		}
	}
	ok, err = f.queue.TopicExist(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed check topic: %w", err)
	}
	if ok {
		err = f.queue.DeleteTopic(ctx, topic)
		if err != nil {
			return fmt.Errorf("failed delete topic: %w", err)
		}
	}

	_, err = f.queue.CreateTopic(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed create topic: %w", err)
	}

	_, err = f.queue.CreateSubscription(ctx, subscription, topic)
	if err != nil {
		return fmt.Errorf("failed create subscription: %w", err)
	}

	return err
}

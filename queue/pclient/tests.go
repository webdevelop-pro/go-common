package pclient

import (
	"context"
	"fmt"

	"github.com/webdevelop-pro/go-common/tests"
)

type contextKey string

const ctxKey contextKey = "db"

func getClient(t tests.TestContext) *Client {
	//nolint:forcetypeassert
	return t.Ctx.Value(pkgName).(*Client)
}

func SendPubSubEvent(topic string, body any, attr map[string]any) tests.SomeAction {
	return func(t tests.TestContext) error {
		_, err := getClient(t).PublishToTopic(t.Ctx, topic, body, attr)
		return err
	}
}

/*
func CheckIfMsgRead(topic string, body any) tests.SomeAction {
	return func(t tests.TestContext) error {
		// _, err := getClient(t).PublishToTopic(context.Background(), topic, body, attr)
		return nil
	}
}
*/

type Fixture struct {
	topic        string
	subscription string
	initData     byte
}

func NewFixture(topic, subscription string, initData byte) Fixture {
	return Fixture{
		topic:        topic,
		subscription: subscription,
		initData:     initData,
	}
}

type FixturesManager struct {
	client   *Client
	fixtures []Fixture
}

func NewFixturesManager(ctx context.Context, fixtures ...Fixture) FixturesManager {
	client, _ := New(ctx)
	return FixturesManager{
		client:   client,
		fixtures: fixtures,
	}
}

func (pubSubF FixturesManager) CleanAndApply() error {
	for _, fixture := range pubSubF.fixtures {
		err := pubSubF.Clean(fixture.topic, fixture.subscription)
		if err != nil {
			return err
		}
	}
	// ToDo
	// Push data to the subscriptions
	// return PubSubF.LoadFixtures(fixtures)
	return nil
}

func (pubSubF FixturesManager) SetCTX(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, pubSubF.client)
}

func (pubSubF FixturesManager) Clean(topic string, subscription string) error {
	ctx := context.Background()
	ok, err := pubSubF.client.SubscriptionExist(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed check subscription: %w", err)
	}
	if ok {
		err := pubSubF.client.DeleteSubscription(ctx, subscription)
		if err != nil {
			return fmt.Errorf("failed delete subscription: %w", err)
		}
	}
	ok, err = pubSubF.client.TopicExist(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed check topic: %w", err)
	}
	if ok {
		err = pubSubF.client.DeleteTopic(ctx, topic)
		if err != nil {
			return fmt.Errorf("failed delete topic: %w", err)
		}
	}

	_, err = pubSubF.client.CreateTopic(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed create topic: %w", err)
	}

	_, err = pubSubF.client.CreateSubscription(ctx, subscription, topic)
	if err != nil {
		return fmt.Errorf("failed create subscription: %w", err)
	}

	return err
}

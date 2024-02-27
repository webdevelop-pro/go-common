package tests

import (
	"context"
	"fmt"

	"github.com/webdevelop-pro/go-common/queue/pclient"
)

type PubSubFixture struct {
	topic        string
	subscription string
	initData     string
}

func NewPubSubFixture(topic, subscription, initData string) PubSubFixture {
	return PubSubFixture{
		topic:        topic,
		subscription: subscription,
		initData:     initData,
	}
}

type PubSubFixturesManager struct {
	client *pclient.Client
}

func NewPubSubFixturesManager(client *pclient.Client) PubSubFixturesManager {
	return PubSubFixturesManager{
		client: client,
	}
}

func (PubSubF PubSubFixturesManager) CleanAndApply(fixtures []PubSubFixture) error {
	for _, fixture := range fixtures {
		err := PubSubF.Clean(fixture.topic, fixture.subscription)
		if err != nil {
			return err
		}
	}
	// ToDo
	// Push data to the subscriptions
	// return PubSubF.LoadFixtures(fixtures)
	return nil

}

func (PubSubF PubSubFixturesManager) Clean(topic string, subscription string) error {
	ctx := context.Background()
	ok, err := PubSubF.client.SubscriptionExist(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed check subscription: %w", err)
	}
	if ok {
		err := PubSubF.client.DeleteSubscription(ctx, subscription)
		if err != nil {
			return fmt.Errorf("failed delete subscription: %w", err)
		}
	}
	ok, err = PubSubF.client.TopicExist(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed check topic: %w", err)
	}
	if ok {
		err = PubSubF.client.DeleteTopic(ctx, topic)
		if err != nil {
			return fmt.Errorf("failed delete topic: %w", err)
		}
	}

	_, err = PubSubF.client.CreateTopic(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed create topic: %w", err)
	}
	PubSubF.client.SetTopic(ctx, topic)

	_, err = PubSubF.client.CreateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed create subscription: %w", err)
	}

	return err
}

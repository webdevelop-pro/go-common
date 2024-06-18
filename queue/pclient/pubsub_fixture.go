package pclient

import (
	"context"
	"fmt"
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
	client *Client
}

func NewPubSubFixturesManager(client *Client) PubSubFixturesManager {
	return PubSubFixturesManager{
		client: client,
	}
}

func (pubSubF PubSubFixturesManager) CleanAndApply(fixtures []PubSubFixture) error {
	for _, fixture := range fixtures {
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

func (pubSubF PubSubFixturesManager) Clean(topic string, subscription string) error {
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

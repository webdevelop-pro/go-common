package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	pclient "github.com/webdevelop-pro/go-common/queue/pclient"
)

var (
	topic        = "test"
	subscription = "test_sub"
)

func TestPublish(t *testing.T) {
	ctx := context.Background()
	pubsubClient, err := pclient.New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	// sub := pubsubClient.client.Subscription(topic)
	// pubsubClient.DeleteSubscription(ctx, "")
	pubsubClient.CreateTopic(ctx, topic)
	t.Run("success publish", func(t *testing.T) {
		msg, err := pubsubClient.Publish(ctx,
			topic,
			map[string]any{"investment_id": 5},
			map[string]string{"ip_address": "31.5.12.199", "request_id": "Xbsdf124d"},
		)
		if err != nil {
			t.Errorf("errors don't match: expected nil, got %s", err)
		}

		id, err := strconv.Atoi(msg.ID)
		if err != nil {
			t.Errorf("pubsub emulator return ID as int, maybe you  are not using it? %s", err)
		}
		if id <= 0 {
			t.Errorf("msg id should be more that 0 %s", err)
		}
	})
}

func TestListenNack(t *testing.T) {
	received_counter := 0

	ctx, cancel := context.WithCancel(context.TODO())
	pubsubClient, err := pclient.New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	pubsubClient.DeleteSubscription(ctx, subscription)
	pubsubClient.CreateTopic(ctx, topic)
	pubsubClient.CreateSubscription(ctx, subscription, topic)

	t.Run("success nack", func(t *testing.T) {
		go pubsubClient.ListenRawMsgs(ctx, subscription, topic, func(ctx context.Context, msg pclient.Message) error {
			received_counter++
			if received_counter%2 != 0 {
				return fmt.Errorf("odd number return an error ... ")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("cannot listen: %s", err)
		}
		pubsubClient.Publish(ctx, topic, map[string]int{"message": 123}, map[string]string{})
		pubsubClient.Publish(ctx, topic, map[string]int{"message": 123}, map[string]string{})
		time.Sleep(time.Second * 4)
		if received_counter != 3 {
			t.Errorf("we should receive same message 3 since listen return an error but got: %d", received_counter)
		}
		cancel()
	})
	pubsubClient.Close()
}

func TestListenAck(t *testing.T) {

	ctx, cancel := context.WithCancel(context.TODO())
	pubsubClient, err := pclient.New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	pubsubClient.DeleteSubscription(ctx, subscription)
	pubsubClient.CreateTopic(ctx, topic)
	pubsubClient.CreateSubscription(ctx, subscription, topic)

	t.Run("success ack", func(t *testing.T) {
		go pubsubClient.ListenEvents(ctx, subscription, topic, func(ctx context.Context, msg pclient.Event) error {
			if msg.ID == "" || msg.Action == "" || msg.ObjectID == 0 || msg.ObjectName == "" {
				return fmt.Errorf("event is empty, its not correct")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("cannot listen: %s", err)
		}
		pubsubClient.Publish(ctx, topic, map[string]int{"message": 123}, map[string]string{})
		pubsubClient.Publish(ctx, topic, map[string]int{"message": 123}, map[string]string{})
		time.Sleep(time.Second * 4)
		cancel()
	})
	pubsubClient.Close()
}

func TestReconnectToNonExistTopic(t *testing.T) {

	ctx, cancel := context.WithCancel(context.TODO())
	pubsubClient, err := pclient.New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	pubsubClient.DeleteSubscription(ctx, subscription)
	// sub.Delete(ctx)
	// pubsubClient.CreateTopic(ctx, topic)
	// pubsubClient.CreateSubscription(ctx, subscription, topic)

	t.Run("success reconnect", func(t *testing.T) {
		go func() {
			go pubsubClient.ListenEvents(ctx, subscription, topic, func(ctx context.Context, msg pclient.Event) error {
				if msg.ID == "" || msg.Action == "" || msg.ObjectID == 0 || msg.ObjectName == "" {
					return fmt.Errorf("event is empty, its not correct")
				}
				return nil
			})
			if err != nil {
				t.Fatalf("cannot connect: %s", err)
			}
		}()
		time.Sleep(time.Second * 5)
		pubsubClient.CreateTopic(ctx, topic)
		pubsubClient.CreateSubscription(ctx, subscription, topic)

		pubsubClient.Publish(ctx, topic, map[string]int{"message": 123}, map[string]string{})
		pubsubClient.Publish(ctx, topic, map[string]int{"message": 123}, map[string]string{})

		time.Sleep(time.Second * 4)
		cancel()
	})
	pubsubClient.Close()
}

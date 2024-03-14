package pclient

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"testing"
	"time"
)

var (
	topic        = "test"
	subscription = "test_sub"
)

func TestPublish(t *testing.T) {
	usr, _ := user.Current()

	os.Setenv("PUBSUB_SERVICE_ACCOUNT_CREDENTIALS", usr.HomeDir+"/.config/gcloud/application_default_credentials.json")
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	ctx := context.Background()
	pubsubClient, err := New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	sub := pubsubClient.client.Subscription(topic)
	sub.Delete(ctx)
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

	usr, _ := user.Current()

	os.Setenv("PUBSUB_SERVICE_ACCOUNT_CREDENTIALS", usr.HomeDir+"/.config/gcloud/application_default_credentials.json")
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	ctx, cancel := context.WithCancel(context.TODO())
	pubsubClient, err := New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	sub := pubsubClient.client.Subscription(topic)
	sub.Delete(ctx)
	pubsubClient.CreateTopic(ctx, topic)
	pubsubClient.CreateSubscription(ctx, subscription, topic)

	t.Run("success nack", func(t *testing.T) {
		err := pubsubClient.ListenRawMsgs(ctx, subscription, topic, func(ctx context.Context, msg Message) error {
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
		if received_counter != 4 {
			t.Errorf("we should receive same message 4 since listen return an error but got: %d", received_counter)
		}
		cancel()
	})
	pubsubClient.Close()
}

func TestListenAck(t *testing.T) {
	usr, _ := user.Current()

	os.Setenv("PUBSUB_SERVICE_ACCOUNT_CREDENTIALS", usr.HomeDir+"/.config/gcloud/application_default_credentials.json")
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	ctx, cancel := context.WithCancel(context.TODO())
	pubsubClient, err := New(ctx)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	sub := pubsubClient.client.Subscription(topic)
	sub.Delete(ctx)
	pubsubClient.CreateTopic(ctx, topic)
	pubsubClient.CreateSubscription(ctx, subscription, topic)

	t.Run("success ack", func(t *testing.T) {
		err := pubsubClient.ListenEvents(ctx, subscription, topic, func(ctx context.Context, msg Event) error {
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

package broker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/webdevelop-pro/go-common/configurator"
)

func TestPubSubPublish(t *testing.T) {
	tests := map[string]struct {
		msg         *Message
		resultID    string
		expectedErr error
	}{
		"success publish": {
			msg:         NewMessage(map[string]any{"investment_id": 5, "ip_address": "31.5.12.199", "request_id": "Xbsdf124d"}, map[string]string{}),
			expectedErr: nil,
		},
	}

	cfg := Config{}
	configurator := configurator.NewConfigurator()
	configurator.New("pubsub", &cfg, "pubsub")
	ctx := context.Background()
	broker, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	broker.CreateTopic(ctx, broker.cfg.Topic)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := broker.Publish(ctx, test.msg)

			if err != test.expectedErr {
				t.Errorf("errors don't match: expected %s, got %s", test.expectedErr, err)
			}
		})
	}
}

func TestPubSubListenNack(t *testing.T) {
	test_msg := NewMessage(map[string]int{"message": 123}, map[string]string{})
	received_counter := 0

	cfg := Config{}
	configurator := configurator.NewConfigurator()
	configurator.New("pubsub", &cfg, "pubsub")
	ctx, cancel := context.WithCancel(context.TODO())
	broker, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	broker.DeleteTopic(ctx, broker.cfg.Topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	sub := broker.client.Subscription(broker.cfg.Topic)
	sub.Delete(ctx)
	broker.CreateTopic(ctx, broker.cfg.Topic)

	t.Run("success nack", func(t *testing.T) {
		err := broker.Listen(ctx, func(ctx context.Context, msg Message) error {
			received_counter++
			if received_counter%2 != 0 {
				return fmt.Errorf("odd ... ")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("cannot listen: %s", err)
		}
		broker.Publish(ctx, test_msg)
		broker.Publish(ctx, test_msg)
		time.Sleep(time.Second * 4)
		if received_counter != 4 {
			t.Errorf("we should receive same message 4 since listen return an error but got: %d", received_counter)
		}
		cancel()
	})
	broker.Close()
}

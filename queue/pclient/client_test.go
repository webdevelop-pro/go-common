package pclient

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/webdevelop-pro/go-common/configurator"
)

func TestPubSubPublish(t *testing.T) {
	msg, err := NewMessage(
		map[string]any{"investment_id": 5},
		map[string]string{"ip_address": "31.5.12.199", "request_id": "Xbsdf124d"},
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tests := map[string]struct {
		msg         *Message
		resultID    string
		expectedErr error
	}{
		"success publish": {
			msg:         msg,
			expectedErr: nil,
		},
	}

	cfg := Config{}
	configurator := configurator.NewConfigurator()
	configurator.New(pkgName, &cfg, pkgName)
	ctx := context.Background()
	pubsubClient, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.CreateTopic(ctx, pubsubClient.cfg.Topic)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := pubsubClient.Publish(ctx, test.msg)

			if err != test.expectedErr {
				t.Errorf("errors don't match: expected %s, got %s", test.expectedErr, err)
			}
		})
	}
}

func TestPubSubListenNack(t *testing.T) {
	test_msg, err := NewMessage(map[string]int{"message": 123}, map[string]string{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	received_counter := 0

	cfg := Config{}
	configurator := configurator.NewConfigurator()
	configurator.New(pkgName, &cfg, pkgName)
	ctx, cancel := context.WithCancel(context.TODO())
	pubsubClient, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("cannot connect %s", err)
	}

	pubsubClient.DeleteTopic(ctx, pubsubClient.cfg.Topic)
	// pubsub emulator has a bug
	// we need to delete subscription manually
	sub := pubsubClient.client.Subscription(pubsubClient.cfg.Topic)
	sub.Delete(ctx)
	pubsubClient.CreateTopic(ctx, pubsubClient.cfg.Topic)
	pubsubClient.CreateSubscription(ctx, pubsubClient.cfg.Subscription)

	t.Run("success nack", func(t *testing.T) {
		err := pubsubClient.Listen(ctx, func(ctx context.Context, msg Message) error {
			received_counter++
			if received_counter%2 != 0 {
				return fmt.Errorf("odd number return an error ... ")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("cannot listen: %s", err)
		}
		pubsubClient.Publish(ctx, test_msg)
		pubsubClient.Publish(ctx, test_msg)
		time.Sleep(time.Second * 4)
		if received_counter != 4 {
			t.Errorf("we should receive same message 4 since listen return an error but got: %d", received_counter)
		}
		cancel()
	})
	pubsubClient.Close()
}

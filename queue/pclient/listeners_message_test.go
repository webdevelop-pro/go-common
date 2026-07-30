package pclient

import (
	"testing"
	"time"

	gpubsub "cloud.google.com/go/pubsub/v2"
)

func TestReceivedMessagePreservesTransportOrderingKey(t *testing.T) {
	t.Parallel()

	attempt := 3
	publishedAt := time.Date(2026, time.July, 30, 12, 0, 0, 0, time.UTC)
	got := receivedMessage(&gpubsub.Message{
		ID:              "message-1",
		Data:            []byte(`{"id":"event-1"}`),
		Attributes:      map[string]string{"type": "investment.created.v1"},
		OrderingKey:     "investment:42",
		PublishTime:     publishedAt,
		DeliveryAttempt: &attempt,
	})

	if got.ID != "message-1" {
		t.Fatalf("ID = %q", got.ID)
	}
	if got.OrderingKey != "investment:42" {
		t.Fatalf("OrderingKey = %q", got.OrderingKey)
	}
	if !got.PublishTime.Equal(publishedAt) {
		t.Fatalf("PublishTime = %s", got.PublishTime)
	}
	if got.Attempt == nil || *got.Attempt != attempt {
		t.Fatalf("Attempt = %v", got.Attempt)
	}
	if string(got.Data) != `{"id":"event-1"}` {
		t.Fatalf("Data = %s", got.Data)
	}
	if got.Attributes["type"] != "investment.created.v1" {
		t.Fatalf("Attributes = %#v", got.Attributes)
	}
}

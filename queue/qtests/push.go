package qtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/webdevelop-pro/go-common/httputils"
	"github.com/webdevelop-pro/go-common/queue/pclient"
	"github.com/webdevelop-pro/go-common/tests"
)

// PushRequest is the envelope Google Pub/Sub sends to push subscribers.
// See https://cloud.google.com/pubsub/docs/push#receive_push.
type PushRequest struct {
	Message      PushMessage `json:"message"`
	Subscription string      `json:"subscription"`
}

type PushMessage struct {
	Attributes      map[string]string `json:"attributes"`
	Data            []byte            `json:"data"`
	MessageID       string            `json:"messageId"`
	PublishTime     time.Time         `json:"publishTime"`
	OrderingKey     string            `json:"orderingKey,omitempty"`
	DeliveryAttempt *int              `json:"deliveryAttempt,omitempty"`
}

// SendPushWebhook emulates a Pub/Sub push delivery of a Webhook to the service
// under test by POSTing a push envelope to its HTTP endpoint.
//
// The service's push handler is expected to invoke its webhook listener
// synchronously and reply 204 only after the work is done; that response is
// the deterministic "done" signal — the test scenario does not need a Sleep.
//
// Host and Port come from the HOST and PORT env vars (the convention shared
// across our integration tests). Use SendPushTo when those env vars don't
// apply.
func SendPushWebhook(path string, msg pclient.Webhook, attrs map[string]string) tests.SomeAction {
	return SendPushTo(os.Getenv("HOST"), os.Getenv("PORT"), path, msg, attrs)
}

// SendPushEvent is the Event-payload counterpart to SendPushWebhook.
func SendPushEvent(path string, msg pclient.Event, attrs map[string]string) tests.SomeAction {
	return SendPushTo(os.Getenv("HOST"), os.Getenv("PORT"), path, msg, attrs)
}

// SendPushTo is the lower-level helper for delivering an arbitrary payload as
// a Pub/Sub push envelope. Use it when the host/port aren't sourced from the
// HOST/PORT env vars or when sending a custom message type.
func SendPushTo(host, port, path string, msg any, attrs map[string]string) tests.SomeAction {
	return func(t tests.TestContext) error {
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal push payload: %w", err)
		}
		envelope := PushRequest{
			Message: PushMessage{
				MessageID:   fmt.Sprintf("test-%d", time.Now().UnixNano()),
				PublishTime: time.Now().UTC(),
				Attributes:  attrs,
				Data:        data,
			},
		}
		body, err := json.Marshal(envelope)
		if err != nil {
			return fmt.Errorf("marshal push envelope: %w", err)
		}
		return tests.SendHTTPRequest(httputils.Request{
			Host:    host,
			Port:    port,
			Method:  http.MethodPost,
			Path:    path,
			Body:    body,
			Headers: map[string]string{"Content-Type": "application/json"},
		}, tests.ExpectedResponse{Code: http.StatusNoContent})(t)
	}
}

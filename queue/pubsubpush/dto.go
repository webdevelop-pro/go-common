// Package pubsubpush provides shared types and HTTP middleware for services
// that receive Pub/Sub messages over a push subscription.
//
// See https://cloud.google.com/pubsub/docs/push#receive_push for the wire
// format and https://cloud.google.com/pubsub/docs/handling-failures for the
// dead-letter / retry semantics.
package pubsubpush

import "time"

// PushRequest is the envelope Google Pub/Sub posts to push subscribers.
//
// DeliveryAttempt is populated at the *top level* of the envelope (not inside
// Message) when the subscription has dead-lettering configured. It tracks how
// many times Pub/Sub has tried to deliver this message.
//
// https://cloud.google.com/pubsub/docs/handling-failures#track_delivery_attempts
type PushRequest struct {
	Message         PushMessage `json:"message"`
	Subscription    string      `json:"subscription"`
	DeliveryAttempt *int        `json:"deliveryAttempt,omitempty"`
}

// PushMessage is the inner Pub/Sub message carried by a push envelope.
type PushMessage struct {
	Attributes  map[string]string `json:"attributes"`
	Data        []byte            `json:"data"`
	MessageID   string            `json:"messageId"`
	PublishTime time.Time         `json:"publishTime"`
	OrderingKey string            `json:"orderingKey,omitempty"`
}

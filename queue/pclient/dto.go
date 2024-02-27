package pclient

import "encoding/json"

type EventType string

const (
	PreAdd     EventType = "pre_add"
	PostAdd    EventType = "post_add"
	PreRemove  EventType = "pre_remove"
	PostRemove EventType = "post_remove"
	PreUpdate  EventType = "pre_update"
	PostUpdate EventType = "post_update"
)

type Message struct {
	Attributes map[string]string
	Data       []byte
	ID         string `json:"message_id"`
}

func NewMessage(data any, attributes map[string]string) (*Message, error) {
	// Returns the message object
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	return &Message{
		Data:       b,
		Attributes: attributes,
	}, nil
}

type Event struct {
	ID         string         `json:"id"`
	Action     EventType      `json:"action" validate:"required"`
	Sender     string         `json:"sender" validate:"required"`
	ObjectID   int            `json:"object_id" validate:"required"`
	ObjectName string         `json:"object_name" validate:"required"`
	RequestID  string         `json:"request_id"`
	IPAddress  string         `json:"ip_address"`
	Data       map[string]any `json:"data"`
}

type Webhook struct {
	ID      string              `json:"id"`
	Action  string              `json:"action" validate:"required"`
	Object  string              `json:"object" validate:"required"`
	Service string              `json:"service" validate:"required"`
	Headers map[string][]string `json:"headers"`
	Data    []byte              `json:"data"`
}

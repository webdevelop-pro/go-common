package broker

import "encoding/json"

type Message struct {
	Attributes map[string]string
	Data       []byte
	ID         string `json:"message_id"`
}

func NewMessage(data any, attributes map[string]string) *Message {
	// Returns the message object
	b, err := json.Marshal(&data)
	if err != nil {
		return nil
	}
	return &Message{
		Data:       b,
		Attributes: attributes,
	}
}

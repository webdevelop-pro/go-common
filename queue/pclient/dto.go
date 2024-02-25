package pclient

import "encoding/json"

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

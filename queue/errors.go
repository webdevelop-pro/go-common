package queue

import "github.com/friendsofgo/errors"

var (
	ErrNotCorrectTopic = errors.New("Topic name incorrect")
	ErrListeningError  = errors.New("ListenWebhooks error")
	ErrAddRoute        = errors.New("Failed add routes")
)

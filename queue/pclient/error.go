package pclient

import "github.com/friendsofgo/errors"

var (
	ErrConfigParse          = errors.New("Cannot parse pubsub config")
	ErrConnection           = errors.New("Connection error")
	ErrCreateSubscription   = errors.New("Cannot create subscription")
	ErrPublish              = errors.New("Cannot publish message")
	ErrNotConnected         = errors.New("Client not connected")
	ErrConnectSubscription  = errors.New("Cannot connect to subscription")
	ErrSubscriptionNotExist = errors.New("Subscription does not exist")
	ErrUnmarshalPubSub      = errors.New("Failed to unmarshal data")
	ErrReceiveSubscription  = errors.New("Failed to receive messages")
	ErrCloseConnection      = errors.New("Close connection error")
	ErrReceiveCallback      = errors.New("Failed to process message")
	ErrTopicNotSet          = errors.New("Cannot connect to empty topic")
	ErrTopicNotExists       = errors.New("Topic does not exist")
	ErrTopicConnect         = errors.New("Cannot connect to topic")
	ErrTopicCreate          = errors.New("Cannot create topic")
)

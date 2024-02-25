package pclient

import "github.com/pkg/errors"

var ErrConfigParse = errors.New("Cannot parse pubsub config")
var ErrConnection = errors.New("Connection error")
var ErrCreateSubscription = errors.New("Cannot create subscription")
var ErrPublish = errors.New("Cannot publish message")
var ErrNotConnected = errors.New("Client not connected")
var ErrConnectSubscription = errors.New("Cannot connect to subscription")
var ErrSubscriptionNotExist = errors.New("Subscription does not exist")
var ErrUnmarshalPubSub = errors.New("Failed to unmarshal message data")
var ErrReceiveSubscription = errors.New("Failed to receive messages")
var ErrCloseConnection = errors.New("Close connection error")
var ErrReceiveCallback = errors.New("Failed to process message")
var ErrTopicNotSet = errors.New("Cannot connect to empty topic")
var ErrTopicNotExists = errors.New("Topic does not exist")
var ErrTopicConnect = errors.New("Cannot connect to topic")
var ErrTopicCreate = errors.New("Cannot create topic")

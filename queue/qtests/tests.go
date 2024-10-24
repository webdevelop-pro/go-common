package qtests

import (
	"context"

	"github.com/webdevelop-pro/go-common/queue/pclient"
	"github.com/webdevelop-pro/go-common/tests"
)

func getQueue(t tests.TestContext) *pclient.Client {
	//nolint:forcetypeassert
	return t.Ctx.Value(queueKey).(*pclient.Client)
}

func SendPubSubEvent(topic string, body any, attr map[string]string) tests.SomeAction {
	return func(t tests.TestContext) error {
		_, err := getQueue(t).PublishToTopic(context.Background(), topic, body, attr)
		return err
	}
}

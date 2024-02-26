package pclient

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/webdevelop-pro/go-common/server/validator"
)

func (b *Client) PublishEvent(ctx context.Context, event Event) (*Message, error) {
	valid := validator.New()
	if err := valid.Verify(event, http.StatusPreconditionFailed); err != nil {
		return nil, err
	}
	attr := map[string]string{}
	return b.PublishToTopic(ctx, b.topic.ID(), event, attr)
}

func (b *Client) PublishWebhook(ctx context.Context, webhook Webhook) (*Message, error) {
	valid := validator.New()
	if err := valid.Verify(webhook, http.StatusPreconditionFailed); err != nil {
		return nil, err
	}
	attr := map[string]string{}
	return b.PublishToTopic(ctx, b.topic.ID(), webhook, attr)
}

func (b *Client) Publish(ctx context.Context, data any, attr map[string]string) (*Message, error) {
	return b.PublishToTopic(ctx, b.topic.ID(), data, attr)
}

func (b *Client) PublishToTopic(ctx context.Context, topicID string, data any, attr map[string]string) (*Message, error) {
	var (
		wg    sync.WaitGroup
		msgID string
		err   error
	)

	t := b.client.Topic(topicID)
	ok, err := t.Exists(ctx)
	if !ok {
		b.log.Error().Err(err).Interface("topic", topicID).Msgf(ErrTopicNotExists.Error())
		return nil, fmt.Errorf("%w: %s", err, b.cfg.Topic)
	}

	msg, err := NewMessage(data, attr)
	if err != nil {
		b.log.Error().Err(err).Interface("data", data).Interface("attr", attr).Msgf(ErrUnmarshalPubSub.Error())
		return nil, err
	}

	wg.Add(1)
	result := t.Publish(ctx, &pubsub.Message{
		Data:       msg.Data,
		Attributes: msg.Attributes,
	})

	go func(res *pubsub.PublishResult) {
		defer wg.Done()
		// The Get method blocks until a server-generated ID or
		// an error is returned for the published message.
		msgID, err = res.Get(ctx)
		if err != nil {
			// Error handling code can be added here.
			b.log.Err(err).Msg(ErrPublish.Error())
			return
		}

		b.log.Debug().Msgf("Published message; msg ID: %v\n", msgID)
	}(result)

	wg.Wait()
	msg.ID = msgID

	return msg, nil
}

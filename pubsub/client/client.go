package client

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/pubsub/broker"
	"github.com/webdevelop-pro/go-logger"
	"google.golang.org/api/option"
)

type Config struct {
	ServiceAccountCredentials string `required:"true" split_words:"true"`
	ProjectID                 string `required:"true" split_words:"true"`
}

type PubsubClient struct {
	client *pubsub.Client
	log    logger.Logger
}

func NewPubsubClient(ctx context.Context) (*PubsubClient, error) {
	log := logger.NewComponentLogger("pubsub-adapter", nil)
	cfg := Config{}

	err := configurator.NewConfiguration(&cfg, "pubsub")
	if err != nil {
		log.Fatal().Err(err).Msg(broker.ErrConfigParse.Error())
	}

	client, err := pubsub.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.ServiceAccountCredentials))
	if err != nil {
		log.Error().Err(err).Msgf(broker.ErrConnection.Error())
		return nil, fmt.Errorf("%w", err)
	}

	return &PubsubClient{
		log:    log,
		client: client,
	}, nil
}

func (p *PubsubClient) PublishMessageToTopic(ctx context.Context, topicID string, attr map[string]string, message []byte) (string, error) {
	var (
		wg    sync.WaitGroup
		msgID string
		err   error
	)

	t := p.client.Topic(topicID)

	wg.Add(1)
	result := t.Publish(ctx, &pubsub.Message{
		Data:       message,
		Attributes: attr,
	})

	go func(res *pubsub.PublishResult) {
		defer wg.Done()
		// The Get method blocks until a server-generated ID or
		// an error is returned for the published message.
		msgID, err = res.Get(ctx)
		if err != nil {
			// Error handling code can be added here.
			p.log.Err(err).Msg(broker.ErrPublish.Error())
			return
		}
		p.log.Info().Msgf("Published message; msg ID: %v\n", msgID)
	}(result)

	wg.Wait()

	return msgID, err
}
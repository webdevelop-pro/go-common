package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/db"
	pubsub "github.com/webdevelop-pro/go-common/pubsub/client"
)

// func SendTestRequest(req *http.Request) ([]byte, int, error) {
// 	httpClient := &http.Client{}
// 	resp, err := httpClient.Do(req)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	bodyBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Errorf("cannot read response body %s", err.Error())
// 		return nil, 0, nil
// 	}

// 	return bodyBytes, resp.StatusCode, nil
// }

type TestContext struct {
	pubsub pubsub.PubsubClient
	db     *db.DB
	t      *testing.T
}

type SomeAction func(t TestContext) (interface{}, error)
type ExpectedResult map[string]interface{}

func SendHttpRequst(method, path string, body []byte) SomeAction {
	return func(t TestContext) (interface{}, error) {
		return false, nil
	}
}

func SendPubSubEvent(topic, body string, attr map[string]string) SomeAction {
	return func(t TestContext) (interface{}, error) {
		msgID, err := t.pubsub.PublishMessageToTopic(context.Background(), topic, attr, []byte(body))

		return msgID, err
	}
}

func SQL(query string, expected ExpectedResult) SomeAction {
	return func(t TestContext) (interface{}, error) {
		var res map[string]interface{}

		query = "select row_to_json(q)::jsonb from (" + query + ") as q"

		err := t.db.QueryRow(context.Background(), query).Scan(&res)
		if err != nil {
			return nil, err
		}

		for key, value := range expected {
			expValue, ok := res[key]
			if assert.True(t.t, ok, fmt.Sprintf("Expected column %s not exist in resukt", key)) {
				assert.Equal(t.t, expValue, value)
			}
		}

		return true, nil
	}
}

func Sleep(d time.Duration) SomeAction {
	return func(t TestContext) (interface{}, error) {
		time.Sleep(d)

		return true, nil
	}
}

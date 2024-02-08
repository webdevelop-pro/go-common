package tests

import (
	"context"
	"time"

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
	pubsubClient pubsub.PubsubClient
}

type SomeAction func(t TestContext) (interface{}, error)

func SendHttpRequst(method, path string, body []byte) SomeAction {
	return func(t TestContext) (interface{}, error) {
		return false, nil
	}
}

func SendPubSubEvent(topic, body string, attr map[string]string) SomeAction {
	return func(t TestContext) (interface{}, error) {
		msgID, err := t.pubsubClient.PublishMessageToTopic(context.Background(), topic, attr, []byte(body))

		return msgID, err
	}
}

func SQL(query string) SomeAction {
	return func(t TestContext) (interface{}, error) {
		return false, nil
	}
}

func Sleep(d time.Duration) SomeAction {
	return func(t TestContext) (interface{}, error) {
		time.Sleep(d)

		return true, nil
	}
}

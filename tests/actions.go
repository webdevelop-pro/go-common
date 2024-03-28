package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/db"
	"github.com/webdevelop-pro/go-common/queue/pclient"
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
	Pubsub pclient.Client
	DB     *db.DB
	T      *testing.T
}

type SomeAction func(t TestContext) error
type ExpectedResult map[string]interface{}

type Request struct {
	Scheme, Host, Method, Path string
	Body                       []byte
	Headers                    map[string]string
}

type ExpectedResponse struct {
	Code int
	Body []byte
}

func SendHttpRequst(req Request, checks ...ExpectedResponse) SomeAction {
	return func(t TestContext) error {
		result, code, err := SendTestRequest(CreateDefaultRequest(req))

		assert.Nil(t.T, err)

		for _, expected := range checks {
			assert.Equal(t.T, expected.Code, code, "Invalid response code")

			if expected.Body != nil {
				CompareJsonBody(t.T, result, expected.Body)
			}
		}

		return nil
	}
}

func SendPubSubEvent(topic string, body any, attr map[string]string) SomeAction {
	return func(t TestContext) error {
		_, err := t.Pubsub.PublishToTopic(context.Background(), topic, body, attr)
		return err
	}
}

func SQL(query string, expected ...ExpectedResult) SomeAction {
	return func(t TestContext) error {
		var res map[string]interface{}

		query = "select row_to_json(q)::jsonb from (" + query + ") as q"

		err := t.DB.QueryRow(context.Background(), query).Scan(&res)
		if err != nil {
			return err
		}

		for _, exp := range expected {
			for key, value := range exp {
				// ToDo:
				// Find library to have colorful compare for maps
				expValue, ok := res[key]
				if assert.True(t.T, ok, fmt.Sprintf("Expected column %s not exist in result", key)) {
          assert.Equal(t.T, value, expValue)
				}
			}
		}

		return nil
	}
}

/*
ToDo:
  - client does not have Listen method
  - Merge broker with client its all client
  - rename pubsub to queue
func CheckPubSubEvent(subID string, msg broker.Message) SomeAction {
	return func(t TestContext) error {
		var equal bool

		log.Trace().Msg("waiting for 5 seconds to get message")

		timer := time.NewTimer(time.Second * 5)
		t.Pubsub.Listen(func(ctx context.Context, msg2 broker.Message) error {
			equal = cmp.Equal(msg, msg2)
			if !equal {
				log.Trace().Interface("msg", msg).Msg("Not equal")
			} else {
				return nil
			}
			return fmt.Errorf("not equal")
		})
		<-timer.C

		return nil
	}
}
*/

// Please usage example
func Sleep(d time.Duration) SomeAction {
	return func(t TestContext) error {
		time.Sleep(d)

		return nil
	}
}

package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
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

type (
	SomeAction     func(t TestContext) error
	ExpectedResult map[string]interface{}
)

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

		query = strings.Replace(query, "\t", " ", -1)
		query = strings.Replace(query, "\n", " ", -1)
		query = strings.Replace(query, "  ", " ", -1)
		row_query := "select row_to_json(q)::jsonb from (" + query + ") as q"

		err := t.DB.QueryRow(context.Background(), row_query).Scan(&res)
		// Do 15 retries automatically
		if err != nil {
			maxRetry := 20
			try := 0
			ticker := time.NewTicker(500 * time.Millisecond)
			for range ticker.C {
				err = t.DB.QueryRow(context.Background(), row_query).Scan(&res)
				if err != nil {
					try++
					if try > maxRetry {
						return errors.Wrapf(err, "for sql: %s", query)
					}
				} else {
					break
				}
			}
		}

		for _, exp := range expected {
			for key, value := range exp {
				// ToDo:
				// Find library to have colorful compare for maps
				expValue, ok := res[key]
				if assert.True(t.T, ok, fmt.Sprintf("Expected column %s not exist in result", key)) {
					value = allowAny(expValue, value)
					assert.Equal(t.T, value, expValue)
				}
			}
		}

		return nil
	}
}

func Sleep(d time.Duration) SomeAction {
	return func(t TestContext) error {
		time.Sleep(d)

		return nil
	}
}

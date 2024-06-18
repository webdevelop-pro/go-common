package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	sqlRetryInterval = 500
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
	// Pubsub pclient.Client
	// DB     *db.DB
	T *testing.T
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

func SendHTTPRequst(req Request, checks ...ExpectedResponse) SomeAction {
	return func(t TestContext) error {
		result, code, err := SendTestRequest(CreateDefaultRequest(req))

		assert.NoError(t.T, err)

		for _, expected := range checks {
			assert.Equal(t.T, expected.Code, code, "Invalid response code")

			if expected.Body != nil {
				CompareJSONBody(t.T, result, expected.Body)
			}
		}

		return nil
	}
}

func SendPubSubEvent(topic string, body any, attr map[string]string) SomeAction {
	return func(t TestContext) error {
		// ToDo: FixME
		/*
			_, err := t.Pubsub.PublishToTopic(context.Background(), topic, body, attr)
			return err
		*/
		return nil
	}
}

func RawSQL(query string) SomeAction {
	return func(t TestContext) error {
		// ToDo: FixME
		/*
			query = strings.ReplaceAll(query, "\t", " ")
			query = strings.ReplaceAll(query, "\n", " ")
			query = strings.ReplaceAll(query, "  ", " ")

			_, err := t.DB.Exec(context.Background(), query)
			return err
		*/
		return nil
	}
}

func SQL(query string, expected ...ExpectedResult) SomeAction {
	return func(t TestContext) error {
		// ToDo: FixME
		/*
			var res map[string]interface{}

			query = strings.ReplaceAll(query, "\t", " ")
			query = strings.ReplaceAll(query, "\n", " ")
			query = strings.ReplaceAll(query, "  ", " ")
			rowQuery := "select row_to_json(q)::jsonb from (" + query + ") as q"

			err := t.DB.QueryRow(context.Background(), rowQuery).Scan(&res)
			// Do 15 retries automatically
			if err != nil {
				maxRetry := 20
				try := 0
				ticker := time.NewTicker(sqlRetryInterval * time.Millisecond)
				for range ticker.C {
					err = t.DB.QueryRow(context.Background(), rowQuery).Scan(&res)
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
						assert.Equal(t.T, expValue, value)
					}
				}
			}
		*/

		return nil
	}
}

func Sleep(d time.Duration) SomeAction {
	return func(_ TestContext) error {
		time.Sleep(d)

		return nil
	}
}

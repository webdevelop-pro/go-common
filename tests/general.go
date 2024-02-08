package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	pubsub "github.com/webdevelop-pro/go-common/pubsub/client"
)

type BodyType string

const (
	JsonBody BodyType = "json"
	FileBody BodyType = "file"
)

type ApiTestCase struct {
	Description      string
	UserID           string
	OnlyForDebugMode bool
	Fixtures         []Fixture

	Request              *http.Request
	BodyType             BodyType
	ExpectedResponseBody []byte
	ExpectedResponseCode int

	TestFunc func(map[string]interface{})
}

type ApiTestCaseV2 struct {
	Description      string
	UserID           string
	OnlyForDebugMode bool

	Fixtures []Fixture
	Actions  []SomeAction
	Checks   []SomeAction
}

func CreateDefaultRequest(method, path string, body []byte) *http.Request {

	appHost := os.Getenv("HOST")
	appPort := os.Getenv("PORT")

	if appHost == "" || appPort == "" {
		log.Fatalf("please set HOST and PORT vars")
	}

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("http://%s:%s%s", appHost, appPort, path),
		bytes.NewBuffer((body)),
	)
	if err != nil {
		log.Fatalf("cannot create new request %s", err.Error())
	}

	req.Header.Add("content-type", "application/json")
	return req
}

func CreateRequestWithFiles(method, path string, body map[string]interface{}, files map[string]string) *http.Request {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	appHost := os.Getenv("HOST")
	appPort := os.Getenv("PORT")

	if appHost == "" || appPort == "" {
		log.Fatalf("please set HOST and PORT vars")
	}

	values := map[string]io.Reader{}
	for k, v := range body {
		values[k] = strings.NewReader(v.(string))
	}

	for k, v := range files {
		f, err := os.Open(v)
		if err != nil {
			log.Fatalf("cannot open file %s", f)
		}
		values[k] = f
	}

	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// upload a file
		if _, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, files[key]); err != nil {
				log.Fatalf("cannot CreateFormFile %s, %s", key, err.Error())
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				log.Fatalf("cannot CreateFormField %s, %s", key, err.Error())
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			log.Fatalf("cannot io.Copy %s, %s", key, err.Error())
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("http://%s:%s%s", appHost, appPort, path),
		buf,
	)
	if err != nil {
		log.Fatalf("cannot create new request %s", err.Error())
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	// Set up content length
	req.Header.Set("Content-Length", fmt.Sprint(req.ContentLength))

	return req
}

func SendTestRequest(req *http.Request) ([]byte, int, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("cannot read response body %s", err.Error())
		return nil, 0, nil
	}

	return bodyBytes, resp.StatusCode, nil
}

func GetPointer(str string) *string {
	return &str
}

func RunApiTest(t *testing.T, Description string, fixtures FixturesManager, scenarios []ApiTestCase) {
	for _, scenario := range scenarios {
		t.Run(scenario.Description, func(t *testing.T) {
			err := fixtures.CleanAndApply(scenario.Fixtures)
			if err != nil {
				assert.Fail(t, "Failed apply fixtures", err)
				log.Panic("Failed apply fixtures")
			}

			result, code, err := SendTestRequest(scenario.Request)

			assert.Nil(t, err)

			assert.Equal(t, scenario.ExpectedResponseCode, code, string(result))

			if scenario.ExpectedResponseBody != nil {
				CompareJsonBody(t, result, scenario.ExpectedResponseBody)
			}

			if scenario.TestFunc != nil {
				bodyMap := make(map[string]interface{})
				if len(result) > 0 {
					if err := json.Unmarshal([]byte(result), &bodyMap); err != nil {
						t.Errorf("cannot convert body %s to map[string]interface, %s", result, err.Error())
					}
				}
				scenario.TestFunc(bodyMap)
			}
		})
	}
}

func RunApiTestV2(t *testing.T, Description string, scenario ApiTestCaseV2) {
	fixtures := NewFixturesManager()
	pubsubClient, _ := pubsub.NewPubsubClient(context.Background())

	t.Run(scenario.Description, func(t *testing.T) {
		testContext := TestContext{
			pubsubClient: *pubsubClient,
		}

		err := fixtures.CleanAndApply(scenario.Fixtures)
		if err != nil {
			assert.Fail(t, "Failed apply fixtures", err)
			log.Panic("Failed apply fixtures")
		}

		for _, action := range scenario.Actions {
			_, err := action(testContext)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, action := range scenario.Checks {
			res, err := action(testContext)
			if err != nil {
				log.Fatal(err)
			}

			if ok, bres := res.(bool); ok {
				assert.True(t, bres)
			}
		}
	})
}

// FixME: use sprew
func CompareJsonBody(t *testing.T, actual, expected []byte) {
	var actualBody, expectedBody map[string]interface{}

	if len(actual) == 0 {
		assert.Fail(t, "server return no data, nothing to compare")
		return
	}

	// Remove tabs and double spaces
	actual = bytes.ReplaceAll(actual, []byte("\t"), []byte(""))
	expected = bytes.ReplaceAll(expected, []byte("\t"), []byte(""))
	actual = bytes.ReplaceAll(actual, []byte("  "), []byte(" "))
	expected = bytes.ReplaceAll(expected, []byte("  "), []byte(" "))

	err := json.Unmarshal(actual, &actualBody)
	if err != nil {
		assert.Failf(t, "failed unmarshal actualBody", "%s, %s", err.Error(), actual)
		return
	}

	err = json.Unmarshal(expected, &expectedBody)
	if err != nil {
		assert.Failf(t, "failed unmarshal expectedBody", "%s %s", err.Error(), expected)
		return
	}

	assert.EqualValuesf(t, expectedBody, actualBody, "responses not equal")
}

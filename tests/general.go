package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/user"
	"reflect"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/db"
	pclient "github.com/webdevelop-pro/go-common/queue/pclient"
)

type BodyType string

const (
	JSONBody BodyType = "json"
	FileBody BodyType = "file"
)

type APITestCase struct {
	Description string
	UserID      string
	// ToDo
	// Write down description why its needed
	OnlyForDebugMode bool
	Fixtures         []Fixture

	Request              *http.Request
	BodyType             BodyType
	ExpectedResponseBody []byte
	ExpectedResponseCode int

	TestFunc func(map[string]interface{})
}

type APITestCaseV2 struct {
	Description string
	UserID      string
	// ToDo
	// Write down description why its needed
	OnlyForDebugMode bool

	Fixtures       []Fixture
	PubSubFixtures []PubSubFixture
	TestActions    []SomeAction
}

// ToDo
// Create in go-common configuration to load env
func LoadEnv(envPath string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get user")
	}

	vars, err := godotenv.Read(envPath)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot read %s", envPath)
	}

	for key, value := range vars {
		value = strings.ReplaceAll(value, "~", usr.HomeDir)
		os.Setenv(key, value)
	}
}

// ToDo
// Create in go-common xserver utils method to make http request
func CreateDefaultRequest(req Request) *http.Request {
	if req.Host == "" {
		appHost := os.Getenv("HOST")
		appPort := os.Getenv("PORT")

		if appHost == "" || appPort == "" {
			log.Fatal().Msg("please set HOST and PORT vars")
		}

		req.Host = appHost + ":" + appPort
	}

	if req.Scheme == "" {
		req.Scheme = "http"
	}

	r, err := http.NewRequest(
		req.Method,
		fmt.Sprintf("%s://%s%s", req.Scheme, req.Host, req.Path),
		bytes.NewBuffer((req.Body)),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot create new request")
	}

	r = r.WithContext(context.Background())
	r.Header.Add("content-type", "application/json")

	for key, value := range req.Headers {
		r.Header.Add(key, value)
	}

	return r
}

// ToDo
// Create in go-common xserver utils method to make request with files
func CreateRequestWithFiles(method, path string, body map[string]interface{}, files map[string]string) *http.Request {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	appHost := os.Getenv("HOST")
	appPort := os.Getenv("PORT")

	if appHost == "" || appPort == "" {
		log.Fatal().Msg("please set HOST and PORT vars")
	}

	values := map[string]io.Reader{}
	for k, v := range body {
		if vs, ok := v.(string); ok {
			values[k] = strings.NewReader(vs)
		}
	}

	for k, v := range files {
		f, err := os.Open(v)
		if err != nil {
			log.Fatal().Err(err).Msgf("cannot open file %s", v)
		}
		values[k] = f
	}

	for key, r := range values {
		var fw io.Writer
		var err error
		x, ok := r.(io.Closer)
		if !ok {
			continue
		}
		// upload a file
		if _, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, files[key]); err != nil {
				w.Close()
				x.Close()
				log.Fatal().Err(err).Msgf("cannot CreateFormFile %s, %s", key, err.Error())
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				w.Close()
				x.Close()
				log.Fatal().Err(err).Msgf("cannot CreateFormField %s, %s", key, err.Error())
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			log.Fatal().Err(err).Msgf("cannot io.Copy %s, %s", key, err.Error())
		}

		x.Close()
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("http://%s%s", net.JoinHostPort(appHost, appPort), path),
		buf,
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot create new request %s", err.Error())
	}

	req = req.WithContext(context.Background())

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
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("cannot read response body %s", err.Error())
		return nil, 0, nil
	}

	return bodyBytes, resp.StatusCode, nil
}

func GetPointer(str string) *string {
	return &str
}

func RunAPITest(t *testing.T, description string, fixtures FixturesManager, scenarios []APITestCase) {
	t.Helper()

	for _, s := range scenarios {
		scenario := s
		t.Run(description+": "+scenario.Description, func(t *testing.T) {
			err := fixtures.CleanAndApply(scenario.Fixtures)
			if err != nil {
				assert.Fail(t, "Failed apply fixtures", err)
				log.Fatal().Err(err).Msgf("Failed apply fixtures")
			}

			result, code, err := SendTestRequest(scenario.Request)

			assert.Nil(t, err)

			assert.Equal(t, scenario.ExpectedResponseCode, code, string(result))

			if scenario.ExpectedResponseBody != nil {
				CompareJSONBody(t, result, scenario.ExpectedResponseBody)
			}

			if scenario.TestFunc != nil {
				bodyMap := make(map[string]interface{})
				if len(result) > 0 {
					if err := json.Unmarshal(result, &bodyMap); err != nil {
						t.Errorf("cannot convert body %s to map[string]interface, %s", result, err.Error())
					}
				}
				scenario.TestFunc(bodyMap)
			}
		})
	}
}

func RunAPITestV2(t *testing.T, description string, scenario APITestCaseV2) {
	t.Helper()

	fixtures := NewFixturesManager()
	pubsubClient, _ := pclient.New(context.Background())
	pubsubFixtures := NewPubSubFixturesManager(&pubsubClient)
	dbClient := db.New()

	t.Run(description+": "+scenario.Description, func(t *testing.T) {
		testContext := TestContext{
			Pubsub: pubsubClient,
			DB:     dbClient,
			T:      t,
		}

		err := fixtures.CleanAndApply(scenario.Fixtures)
		if err != nil {
			assert.Fail(t, "Failed apply fixtures", err)
			log.Fatal().Err(err).Msgf("Failed apply fixtures")
		}
		err = pubsubFixtures.CleanAndApply(scenario.PubSubFixtures)
		if err != nil {
			assert.Fail(t, "Failed apply pubsub fixtures", err)
			log.Fatal().Err(err).Msgf("Failed apply pubsub fixtures")
		}

		for _, action := range scenario.TestActions {
			err := action(testContext)
			if err != nil {
				log.Fatal().Err(err).Msgf("scenario return an error")
			}
		}
	})
}

func allowDictAny(src, dst map[string]interface{}) map[string]interface{} {
	res := dst

	for k, v := range dst {
		switch val := v.(type) {
		case string:
			if val == "%any%" && src != nil && !reflect.ValueOf(src[k]).IsZero() {
				res[k] = src[k]
			}
		case int:
			if val == math.MinInt {
				dst[k] = src[k]
			}
		case map[string]any:
			if srck, ok := src[k].(map[string]any); ok {
				res[k] = allowDictAny(srck, val)
			}
		}
	}

	return res
}

func allowAny(src, dst interface{}) interface{} {
	res := dst

	switch val := dst.(type) {
	case string:
		if val == "%any%" && src != nil && !reflect.ValueOf(src).IsZero() {
			res = src
		}
	case int:
		if val == math.MinInt {
			res = src
		}
	}

	return res
}

// ToDo
// use sprew or other library to better show different in maps
func CompareJSONBody(t *testing.T, actual, expected []byte) {
	t.Helper()

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

	expectedBody = allowDictAny(actualBody, expectedBody)
	assert.EqualValuesf(t, expectedBody, actualBody, "responses not equal")
}

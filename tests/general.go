package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	Fixture          string

	Method, URL          string
	Body                 io.Reader
	BodyType             BodyType
	ExpectedResponseBody *string
	ExpectedResponseCode int

	TestFunc func()
}

func SendTestRequest(method, path, userID string, body io.Reader) (string, int, error) {
	httpClient := &http.Client{}

	appHost := os.Getenv("TESTS_APP_HOST")
	appPort := os.Getenv("TESTS_APP_PORT")

	if appHost == "" || appPort == "" {
		return "", 0, fmt.Errorf("please set TESTS_APP host and port vars")
	}

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%s%s", appHost, appPort, path), body)
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("content-type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	bodyString := ""
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", 0, nil
		}
		bodyString = string(bodyBytes)
	}

	return bodyString, resp.StatusCode, nil
}

func GetPointer(str string) *string {
	return &str
}

func RunApiTest(t *testing.T, fixtures FixturesManager, scenario ApiTestCase) {
	t.Run(scenario.Description, func(t *testing.T) {
		defer func() {
			err := fixtures.Clean()
			if err != nil {
				assert.Fail(t, "Failed clean fixtures", err)
			}
		}()

		err := fixtures.CleanAndApply(scenario.Fixture)
		if err != nil {
			assert.Fail(t, "Failed apply fixtures", err)
		}

		if scenario.Method == "" {
			scenario.Method = "GET"
		}

		result, code, err := SendTestRequest(scenario.Method, scenario.URL, scenario.UserID, scenario.Body)

		assert.Nil(t, err)

		assert.Equal(t, scenario.ExpectedResponseCode, code)

		if scenario.ExpectedResponseBody != nil {
			if scenario.BodyType == FileBody {
				CompareFilesBody(t, result, *scenario.ExpectedResponseBody)
			} else {
				CompareJsonBody(t, result, *scenario.ExpectedResponseBody)
			}
		}

		if scenario.TestFunc != nil {
			scenario.TestFunc()
		}
	})
}

// FixME: use sprew
func CompareFilesBody(t *testing.T, actual, expected string) {
	expected = strings.ReplaceAll(expected, "\t", "")
	assert.Equal(t, expected, actual)
}

// FixME: use sprew
func CompareJsonBody(t *testing.T, actual, expected string) {
	var actualBody, expectedBody map[string]interface{}

	err := json.Unmarshal([]byte(actual), &actualBody)
	if err != nil {
		assert.Fail(t, "failed unmarshal actualBody")
		return
	}

	err = json.Unmarshal([]byte(expected), &expectedBody)
	if err != nil {
		assert.Fail(t, "failed unmarshal expectedBody")
		return
	}

	checkMaps(t, actualBody, expectedBody, "")
}

// FixME: use sprew
func checkMaps(t *testing.T, dict, subDict map[string]interface{}, baseField string) {
	for key, subValue := range subDict {
		fieldPath := fmt.Sprintf("%s_%s", baseField, key)

		value, ok := dict[key]
		if !ok {
			assert.FailNow(t, fieldPath+" not exist")
			return
		}

		if subValue == nil {
			if value == nil {
				continue
			} else {
				assert.Fail(t, fieldPath+" not nil")
			}
		}

		if reflect.TypeOf(value) != reflect.TypeOf(subValue) {
			assert.Fail(t, fieldPath+" different types")
			return
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Array, reflect.Slice:
			checkArrays(t, value.([]interface{}), subValue.([]interface{}), fieldPath)
		case reflect.Map:
			checkMaps(t, value.(map[string]interface{}), subValue.(map[string]interface{}), fieldPath)
		default:
			assert.Equalf(t, subValue, value, fieldPath)
		}
	}
}

// FixME: use sprew
func checkArrays(t *testing.T, array, subArray []interface{}, baseField string) {
	if len(array) != len(subArray) {
		assert.FailNow(t, baseField+" expected array and actual have different length")
		return
	}

	for index, subValue := range subArray {
		fieldPath := fmt.Sprintf("%s_%d", baseField, index)

		value := array[index]

		if reflect.TypeOf(value) != reflect.TypeOf(subValue) {
			assert.Fail(t, fieldPath+" different types")
			return
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Array, reflect.Slice:
			checkArrays(t, value.([]interface{}), subValue.([]interface{}), fieldPath)
		case reflect.Map:
			checkMaps(t, value.(map[string]interface{}), subValue.(map[string]interface{}), fieldPath)
		default:
			assert.Equalf(t, subValue, value, fieldPath)
		}
	}
}

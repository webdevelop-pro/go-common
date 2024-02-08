//nolint:paralleltest,thelper
package test

import (
	"os"
	"testing"
	"time"

	. "github.com/webdevelop-pro/go-common/tests"
)

func TestMain(m *testing.M) {
	LoadEnv(".env.tests")

	// go start.Server()

	os.Exit(m.Run())
}

func TestExample(t *testing.T) {
	RunApiTestV2(t,
		"",
		ApiTestCaseV2{
			Description: "Example case",
			Fixtures:    []Fixture{},
			Actions: []SomeAction{
				// SendHttpRequst("POST", "/events/sendgrid/test_topic?object=email&action=update&auth_type=auto&auth_token=XXXXX", []byte(`{"test": "message"}`)),
				SendPubSubEvent("test_topic", "{}", map[string]string{}),
			},
			Checks: []SomeAction{
				Sleep(time.Second * 2),
				SQL(
					"select 1 as col_1, 'a' as col_2 limit 1",
					ExpectedResult{
						"col_1": 1.0,
						"col_2": "a",
					},
				),
			},
		},
	)
}

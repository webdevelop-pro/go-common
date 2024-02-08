//nolint:paralleltest,thelper
package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/webdevelop-pro/go-common/tests"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(".env.tests")
	if err != nil {
		log.Fatalln(err)
	}

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
			},
		},
	)
}

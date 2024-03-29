//nolint:paralleltest,thelper
package test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webdevelop-pro/go-common/queue/pclient"
	. "github.com/webdevelop-pro/go-common/tests"
)

/*
PROBLEM:

- this line set up INCORRECT envs for me
- so I cannot run test with it

func TestMain(m *testing.M) {
	LoadEnv(".env.tests")

	// go start.Server()

	os.Exit(m.Run())
}
*/

func TestExample(t *testing.T) {
	RunApiTestV2(t,
		"",
		ApiTestCaseV2{
			Description: "HTTP & PubSub example",
			Fixtures:    []Fixture{},
			PubSubFixtures: []PubSubFixture{
				NewPubSubFixture(
					os.Getenv("PUBSUB_TOPIC"),
					os.Getenv("PUBSUB_SUBSCRIPTION"),
					"",
				),
			},
			TestActions: []SomeAction{
				// SendHttpRequst("POST", "/events/sendgrid/test_topic?object=email&action=update&auth_type=auto&auth_token=XXXXX", []byte(`{"test": "message"}`)),
				SendPubSubEvent(os.Getenv("PUBSUB_TOPIC"), "{}", map[string]string{}),
				Sleep(time.Second * 2),
				SQL(
					"select 1 as col_1, 'a' as col_2 limit 1",
					ExpectedResult{
						"col_1": 1.0,
						"col_2": "a",
					},
				),
				SendPubSubEvent(
					os.Getenv("PUBSUB_TOPIC"),
					pclient.Webhook{
						Object:  "profile",
						Action:  "update_accr",
						Service: "north_capital",
						Data: []byte(`
								"accountId":["NO_INVESTMENTS"],
								"airequestId":["Tzboaa"],
								"aiRequestStatus":["Approved"],
								"accreditedStatus":["Verified Accredited"],
							}`),
					},
					map[string]string{},
				),
				SendHttpRequst(
					Request{
						Scheme: "https",
						Host:   "google.com",
						Headers: map[string]string{
							"identity_id": "test",
						},
						Method: "GET",
						Path:   "/",
					},
					ExpectedResponse{
						Code: 200,
					},
				),
				func(tc TestContext) error {
					// Some custom test code
					// tc.DB.Query ...
					// tc.Pubsub....

					assert.True(tc.T, true)

					return nil
				},
			},
		},
	)
}

// TestAny word %any% allow to have any string
func TestAny(t *testing.T) {
	CompareJsonBody(
		t,
		[]byte(`{"nc_issuer_id":"10010949", "nc_issuer_status":"Pending"}`),
		[]byte(`{"nc_issuer_id":"%any%", "nc_issuer_status":"Pending"}`),
	)
}

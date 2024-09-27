//nolint:paralleltest,thelper
package qtests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/webdevelop-pro/go-common/queue/pclient"
	gTests "github.com/webdevelop-pro/go-common/tests"
)

func TestExample(t *testing.T) {
	ctx := context.TODO()
	gTests.RunTableTest(t, ctx,
		[]gTests.FixturesManager{
			NewFixturesManager(ctx, NewFixture(os.Getenv("PUBSUB_TOPIC_WEBHOOK"), os.Getenv("PUBSUB_TOPIC_WEBHOOK"), "")),
		},
		gTests.TableTest{
			Description: "test of the table test",
			Scenarios: []gTests.TestScenario{
				{
					Description: "Success test",
					TestActions: []gTests.SomeAction{
						// SendHttpRequst("POST", "/events/sendgrid/test_topic?object=email&action=update&auth_type=auto&auth_token=XXXXX", []byte(`{"test": "message"}`)),
						SendPubSubEvent(os.Getenv("PUBSUB_TOPIC_WEBHOOK"), "{}", map[string]string{}),
						gTests.Sleep(time.Second * 2),
						SendPubSubEvent(
							os.Getenv("PUBSUB_TOPIC_WEBHOOK"),
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
					},
				},
			},
		},
	)
}

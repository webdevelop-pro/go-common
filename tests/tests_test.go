//nolint:paralleltest,thelper
package tests

/*
func TestExample(t *testing.T) {
	ctx := context.TODO()
	pubsub, err := pclient.New(ctx)
	if err != nil {
		assert.Fail(t, "Failed to create pubsub client", err)
	}
	pubsub.CreateTopic(ctx, os.Getenv("PUBSUB_TOPIC_WEBHOOK"))
	pubsub.CreateSubscription(ctx, os.Getenv("PUBSUB_TOPIC_WEBHOOK"), os.Getenv("PUBSUB_SUBSCRIPTION_WEBHOOK"))
	RunAPITestV2(t,
		"",
		APITestCaseV2{
			Description: "HTTP & PubSub example",
			Fixtures:    []Fixture{},
			PubSubFixtures: []PubSubFixture{
				NewPubSubFixture(
					os.Getenv("PUBSUB_TOPIC_WEBHOOK"),
					os.Getenv("PUBSUB_SUBSCRIPTION_WEBHOOK"),
					"",
				),
			},
			TestActions: []SomeAction{
				// SendHttpRequst("POST", "/events/sendgrid/test_topic?object=email&action=update&auth_type=auto&auth_token=XXXXX", []byte(`{"test": "message"}`)),
				SendPubSubEvent(os.Getenv("PUBSUB_TOPIC_WEBHOOK"), "{}", map[string]string{}),
				Sleep(time.Second * 2),
				SQL(
					"select 1 as col_1, 'a' as col_2 limit 1",
					ExpectedResult{
						"col_1": 1.0,
						"col_2": "a",
					},
				),
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
				SendHTTPRequst(
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
	CompareJSONBody(
		t,
		[]byte(`{"nc_issuer_id":"10010949", "nc_issuer_status":"Pending"}`),
		[]byte(`{"nc_issuer_id":"%any%", "nc_issuer_status":"Pending"}`),
	)
}
*/

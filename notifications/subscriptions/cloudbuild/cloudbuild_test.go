// Package p contains a Pub/Sub Cloud Function.

package cloudbuild

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessEvent(t *testing.T) {
	t.Parallel()

	worker := NewWorker(
		Config{
			SlackToken:        "xoxb-2523150300577-2495865134151-7yze2ClHQ1zLGILNvpGOYJOm",
			GithubAccessToken: "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",
			GitRepoOwner:      "replier-ai",
			Channels: map[string][]Channel{
				"all": {
					{
						Type: Slack,
						To:   "tests",
					},
					{
						Type: Matrix,
						To:   "https://matrix.unfederalreserve.com/api/v1/matrix/hook/Z0IGHMnoLWRU8WDPx26FAwIdjJmbMaVX2Duy6LYQlmqUYTwPe77ZnUh4aKNhSXxv",
					},
				},
			},
		},
	)

	err := worker.ProcessEvent(context.Background(), PubSubMessage{
		Data: []byte(`
			{
				"id": "ee57301f-580f-4563-b2c2-ae6f137fb91a",
				"status": "SUCCESS",
				"startTime": "2021-09-19T10:47:24.334389196Z",
				"finishTime": "2021-09-19T10:48:26.652347Z",
				"projectId": "replier-ai-dev",
				"logUrl": "https://console.cloud.google.com/cloud-build/builds/ee57301f-580f-4563-b2c2-ae6f137fb91a?project=639894471322",
				"substitutions": {
					"REVISION_ID": "9e1180c77d205fbfa9052bcfd493928c430033f0",
					"SHORT_SHA": "9e1180c",
					"BRANCH_NAME": "cloudbuild",
					"REPO_NAME": "payment-api",
					"COMMIT_SHA": "9e1180c77d205fbfa9052bcfd493928c430033f0",
					"_SERVICE_NAME": "payment-api",
					"REF_NAME": "cloudbuild",
					"TRIGGER_NAME": "pament-api",
					"TRIGGER_BUILD_CONFIG_PATH": "cloudbuild.yaml"
				}
			}`,
		),
	})

	require.Nil(t, err)
}

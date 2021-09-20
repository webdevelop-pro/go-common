// Package p contains a Pub/Sub Cloud Function.
package cloudbuild

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/webdevelop-pro/go-common/git"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/notifications/senders"
)

type ChannelType string

const (
	Slack  ChannelType = "slack"
	Matrix ChannelType = "matrix"
)

type ChannelsMap map[string][]Channel

type Channel struct {
	Type ChannelType `json:"type"`
	To   string      `json:"to"`
}

func (channels *ChannelsMap) Decode(value string) error {
	return json.Unmarshal([]byte(value), channels)
}

type Config struct {
	SlackToken        string      `required:"true" split_words:"true"`
	Channels          ChannelsMap `required:"true" split_words:"true"`
	GitRepoOwner      string      `required:"true" split_words:"true"`
	GithubAccessToken string      `required:"true" split_words:"true"`
}

type Worker struct {
	Config

	gitClient git.Client
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type EventRecord struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	LogURL        string    `json:"logUrl"`
	StartTime     time.Time `json:"startTime"`
	FinishTime    time.Time `json:"finishTime"`
	ProjectId     string    `json:"projectId"`
	Substitutions struct {
		RepoName   string `json:"REPO_NAME"`
		CommitSha  string `json:"COMMIT_SHA"`
		ShortSha   string `json:"SHORT_SHA"`
		BranchName string `json:"BRANCH_NAME"`
	} `json:"substitutions"`
}

const baselink = "https://console.cloud.google.com/logs/query;query=resource.type%3D%22k8s_container%22%0Atimestamp%3D%22"

// Subscribe consumes a Pub/Sub message.
func Subscribe(ctx context.Context, m PubSubMessage) error {
	var conf Config
	log := logger.GetDefaultLogger(nil)

	err := envconfig.Process("", &conf)
	if err != nil {
		log.Error().Err(err).Msg("invalid config")
		return err
	}

	worker := NewWorker(conf)

	err = worker.ProcessEvent(ctx, m)
	if err != nil {
		log.Error().Err(err).Msg("invalid config")
	}

	return err
}

func NewWorker(conf Config) Worker {
	return Worker{
		Config: conf,
		gitClient: git.GithubClient{
			AccessToken: conf.GithubAccessToken,
			RepoOwner:   conf.GitRepoOwner,
		},
	}
}

func (w Worker) ProcessEvent(ctx context.Context, m PubSubMessage) error {
	var event EventRecord
	json.Unmarshal(m.Data, &event)

	if event.Status != "SUCCESS" && event.Status != "FAILURE" && event.Status != "TIMEOUT" {
		return nil
	}

	send := func(output []Channel) error {
		for _, channel := range output {
			err := w.sendNotification(event, channel)
			if err != nil {
				return err
			}
		}

		return nil
	}

	output, ok := w.Channels[event.Substitutions.RepoName]
	if ok {
		err := send(output)
		if err != nil {
			return err
		}
	}

	err := send(w.Channels["all"])
	if err != nil {
		return err
	}

	return nil
}

func (w Worker) sendNotification(event EventRecord, channel Channel) error {
	message, err := w.createMessage(event, channel.Type)
	if err != nil {
		return fmt.Errorf("Failed create notification message: %w", err)
	}

	routes := map[ChannelType]senders.Send{
		Matrix: senders.SendToMatrix,
		Slack:  senders.SlackSender{Token: w.SlackToken}.SendToSlack,
	}

	err = routes[channel.Type](
		message,
		channel.To,
		senders.MessageStatus(event.Status),
	)
	if err != nil {
		return fmt.Errorf("Failed send notification: %w", err)
	}

	return nil
}

func (w Worker) createMessage(event EventRecord, channelType ChannelType) (string, error) {
	duration := event.FinishTime.Sub(event.StartTime)

	commit, err := w.gitClient.GetCommit(event.Substitutions.RepoName, event.Substitutions.CommitSha)
	if err != nil {
		return "", err
	}

	messageTemplates := map[ChannelType]string{
		Slack:  "Build <%s|%s>, <https://github.com/%s/%s/commit/%s|%s/%s - %s>, Duration: %s\n%s: %s",
		Matrix: "<p>Build <a href=\"%s\">%s</a>, <a href=\"https://github.com/%s/%s/commit/%s\">%s/%s - %s</a>, %s: %s, Duration: %s</p>",
	}

	return fmt.Sprintf(
		messageTemplates[channelType],
		event.LogURL,
		event.Status,
		event.ProjectId,
		event.Substitutions.RepoName,
		event.Substitutions.CommitSha,
		event.Substitutions.RepoName,
		event.Substitutions.BranchName,
		event.Substitutions.ShortSha,
		humanizeDuration(duration),
		commit.Author,
		commit.Message,
	), nil
}

func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d sec", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d min %d sec", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d min %d sec",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d min %d sec",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}

package git

import (
	"context"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type Client interface {
	getCommit(repo, sha string) (GitCommit, error)
}

type GitCommit struct {
	Author  string
	URL     string
	Message string
}

type GithubClient struct {
	AccessToken string `json:"access_token"`
	RepoOwner   string `json:"repo_owner"`
}

func (g GithubClient) getCommit(repo, sha string) (GitCommit, error) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	commit, _, err := client.Git.GetCommit(
		ctx,
		g.RepoOwner,
		repo,
		sha,
	)

	return GitCommit{
		Author:  *commit.Author.Name,
		URL:     *commit.URL,
		Message: *commit.Message,
	}, err
}

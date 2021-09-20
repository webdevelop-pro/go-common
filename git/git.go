package git

import (
	"context"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// Client ...
type Client interface {
	GetCommit(repo, sha string) (*Commit, error)
}

// Commit ...
type Commit struct {
	Author  string
	URL     string
	Message string
}

// GithubClient ...
type GithubClient struct {
	AccessToken string `json:"access_token"`
	RepoOwner   string `json:"repo_owner"`
}

// GetCommit ...
func (g GithubClient) GetCommit(repo, sha string) (*Commit, error) {
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
	if err != nil {
		return nil, err
	}

	return &Commit{
		Author:  *commit.Author.Name,
		URL:     *commit.URL,
		Message: *commit.Message,
	}, err
}

package util

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"regexp"
	"strings"
)

type GithubClient struct {
	owner      string
	Repository string
	client     *github.Client
	ctx        context.Context
}

func NewGithubClient() *GithubClient {
	ctx := context.Background()
	ow, repo := currentGitRepository()

	return &GithubClient{
		owner:      ow,
		Repository: repo,
		ctx:        ctx,
		client:     newInnerClient(ctx),
	}
}

func newInnerClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("github_access_token")},
	)

	c := oauth2.NewClient(ctx, ts)

	return github.NewClient(c)
}

func currentGitRepository() (string, string) {
	url := ConfigGet("remote.origin.url")

	//url should be git@github.com:kunicmarko20/deko-cli.git
	//now we need to remove prefix and suffix to get owner and repository
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")

	return parts[0], parts[1]
}

func (gc *GithubClient) CreateReleasePullRequest(pb ProgressBar, rlsDate string, rlsBranch string) *github.PullRequest {
	t := pb.Start("creating release pull request")

	prTitle := "Release " + rlsDate
	bsBr := "master"
	mtCanModify := true

	pr, _, err := gc.client.PullRequests.Create(
		gc.ctx,
		gc.owner,
		gc.Repository,
		&github.NewPullRequest{
			Title:               &prTitle,
			Head:                &rlsBranch,
			Base:                &bsBr,
			MaintainerCanModify: &mtCanModify,
		})

	if err != nil {
		Exit(err)
	}

	t.Increment(1)

	return pr
}

func (gc *GithubClient) UpdateReleasePullRequestBodyWithTickets(pb ProgressBar, prNum int) {
	t := pb.Start("updating release pull request body")

	_, _, err := gc.client.PullRequests.Edit(
		gc.ctx,
		gc.owner,
		gc.Repository,
		prNum,
		&github.PullRequest{
			Body: gc.newPullRequestBodyFromCommits(prNum),
		})

	if err != nil {
		Exit(err)
	}

	t.Increment(1)
}

func (gc *GithubClient) newPullRequestBodyFromCommits(prNum int) *string {
	cmts, _, err := gc.client.PullRequests.ListCommits(gc.ctx, gc.owner, gc.Repository, prNum, nil)
	if err != nil {
		Exit(err)
	}

	var newPRBody string
	reg, _ := regexp.Compile("\\[\\w+\\-\\d+\\]")
	for _, cmt := range cmts {
		if reg.MatchString(*cmt.Commit.Message) {
			newPRBody += *cmt.Commit.Message + "\n"
		}
	}

	return &newPRBody
}

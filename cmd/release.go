package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/ldez/go-git-cmd-wrapper/config"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
	"time"
)

var (
	releaseCmd = &cobra.Command{
		Use: "release",
		Run: NewReleaseCommand().Run,
	}
)

type ReleaseCommand struct {
	owner string
	repository string
	client *github.Client
	ctx context.Context
}

func NewReleaseCommand() *ReleaseCommand {
	ow, repo := currentGitRepository()

	return &ReleaseCommand{
		owner: ow,
		repository: repo,
		client: github.NewClient(nil),
		ctx: context.Background(),
	}
}

func currentGitRepository() (string, string) {
	url, err := git.Config(config.Get("remote.origin.url", ""))
	if err != nil {
		er(err)
	}

	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")

	return parts[0], parts[1]
}

func (r *ReleaseCommand) Run(cmd *cobra.Command, args []string) {
	rlsDate := time.Now().Format("20200706")
	rlsBranch := "release-" + rlsDate

	fmt.Print("creating release branch")
	r.createReleaseBranch(rlsBranch)

	fmt.Print("creating release pull request")
	prNum := r.createReleasePullRequest(rlsDate, rlsBranch)

	fmt.Print("updating release pull request body")
	r.updateReleasePullRequestBodyWithTickets(*prNum)

}

func (r *ReleaseCommand) createReleaseBranch(rlsBranch string) {
	ref := "refs/heads/" + rlsBranch
	_, _, err := r.client.Git.CreateRef(
		r.ctx,
		r.owner,
		r.repository,
		&github.Reference{
			Ref: &ref,
			Object: &github.GitObject{SHA: r.newestSha()},
		})

	if err != nil {
		er(err)
	}
}

func (r *ReleaseCommand) newestSha() *string {
	br, _, err := r.client.Repositories.GetBranch(r.ctx, r.owner, r.repository, "branch")

	if err != nil {
		er(err)
	}

	return br.Commit.SHA
}

func (r *ReleaseCommand) createReleasePullRequest(rlsDate string, rlsBranch string) *int {
	prTitle := "Release" + rlsDate
	bsBr := "master"
	mtCanModify := true

	pr, _, err := r.client.PullRequests.Create(
		r.ctx,
		r.owner,
		r.repository,
		&github.NewPullRequest{
			Title: &prTitle,
			Head: &rlsBranch,
			Base: &bsBr,
			MaintainerCanModify: &mtCanModify,
		})

	if err != nil {
		er(err)
	}

	return pr.Number
}

func (r *ReleaseCommand) updateReleasePullRequestBodyWithTickets(prNum int) {
	cmts, _, err := r.client.PullRequests.ListCommits(r.ctx, r.owner, r.repository, prNum, nil)
	if err != nil {
		er(err)
	}

	for _, cmt := range cmts {
		r := regexp.MustCompile(`\[\w+\-\d+\]`)
		r.FindStringSubmatch(*cmt.Commit.Message)
	}
}
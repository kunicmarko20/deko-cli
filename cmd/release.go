package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jedib0t/go-pretty/progress"
	"github.com/ldez/go-git-cmd-wrapper/checkout"
	"github.com/ldez/go-git-cmd-wrapper/config"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/merge"
	"github.com/ldez/go-git-cmd-wrapper/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
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
	owner      string
	repository string
	client     *github.Client
	ctx        context.Context
	pw         progress.Writer
}

func NewReleaseCommand() *ReleaseCommand {
	ctx := context.Background()
	ow, repo := currentGitRepository()

	pw := progress.NewWriter()
	pw.SetTrackerLength(10)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.SetStyle(progress.StyleBlocks)
	pw.SetMessageWidth(50)
	pw.SetNumTrackersExpected(7)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 50)
	pw.Style().Colors = progress.StyleColorsExample

	return &ReleaseCommand{
		owner:      ow,
		repository: repo,
		ctx:        ctx,
		pw:         pw,
	}
}

func currentGitRepository() (string, string) {
	url, err := git.Config(config.Get("remote.origin.url", ""))
	if err != nil {
		er(err)
	}

	//url should be git@github.com:kunicmarko20/deko-cli.git
	//now we need to remove prefix and suffix to get owner and repository
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "git@github.com:")
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")

	return parts[0], parts[1]
}

func (r *ReleaseCommand) Run(cmd *cobra.Command, args []string) {
	r.client = githubClient(r.ctx)

	now := time.Now()
	rlsDate := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	rlsBranch := "release-" + rlsDate

	go r.pw.Render()
	r.createReleaseBranch(rlsBranch)

	tracker1 := progress.Tracker{Message: "creating release pull request", Total: 1}
	r.pw.AppendTracker(&tracker1)
	pr := r.createReleasePullRequest(rlsDate, rlsBranch)
	tracker1.Increment(1)

	tracker2 := progress.Tracker{Message: "updating release pull request body", Total: 1}
	r.pw.AppendTracker(&tracker2)
	r.updateReleasePullRequestBodyWithTickets(*pr.Number)
	tracker2.Increment(1)

	for r.pw.IsRenderInProgress() {
		if r.pw.LengthActive() == 0 {
			r.pw.Stop()
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println(fmt.Sprintf("done, visit %s to see your pull request", *pr.HTMLURL))
	err := webbrowser.Open(*pr.HTMLURL)
	if err != nil {
		er(err)
	}
}

func githubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("github_access_token")},
	)

	c := oauth2.NewClient(ctx, ts)

	return github.NewClient(c)
}

func (r *ReleaseCommand) createReleaseBranch(rlsBranch string) {
	tracker1 := progress.Tracker{Message: "changing branch to master", Total: 1}
	r.pw.AppendTracker(&tracker1)
	msg, err := git.Checkout(checkout.Branch("master"))
	if err != nil {
		er(msg)
	}
	tracker1.Increment(1)

	tracker2 := progress.Tracker{Message: "pulling new changes", Total: 1}
	r.pw.AppendTracker(&tracker2)
	msg, err = git.Pull()
	if err != nil {
		er(msg)
	}
	tracker2.Increment(1)

	tracker3 := progress.Tracker{Message: fmt.Sprintf("creating new branch '%s'", rlsBranch), Total: 1}
	r.pw.AppendTracker(&tracker3)
	msg, err = git.Checkout(checkout.NewBranch(rlsBranch))
	if err != nil {
		er(msg)
	}
	tracker3.Increment(1)

	tracker4 := progress.Tracker{Message: "merging origin/staging", Total: 1}
	r.pw.AppendTracker(&tracker4)
	msg, err = git.Merge(merge.Commits("origin/staging"))
	if err != nil {
		er(msg)
	}
	tracker4.Increment(1)

	tracker5 := progress.Tracker{Message: fmt.Sprintf("pushing new branch '%s' to origin", rlsBranch), Total: 1}
	r.pw.AppendTracker(&tracker5)
	msg, err = git.Push(push.Remote("origin"), push.Remote(rlsBranch))
	if err != nil {
		er(msg)
	}
	tracker5.Increment(1)
}

func (r *ReleaseCommand) newestSha() *string {
	br, _, err := r.client.Repositories.GetBranch(r.ctx, r.owner, r.repository, "branch")

	if err != nil {
		er(err)
	}

	return br.Commit.SHA
}

func (r *ReleaseCommand) createReleasePullRequest(rlsDate string, rlsBranch string) *github.PullRequest {
	prTitle := "Release " + rlsDate
	bsBr := "master"
	mtCanModify := true

	pr, _, err := r.client.PullRequests.Create(
		r.ctx,
		r.owner,
		r.repository,
		&github.NewPullRequest{
			Title:               &prTitle,
			Head:                &rlsBranch,
			Base:                &bsBr,
			MaintainerCanModify: &mtCanModify,
		})

	if err != nil {
		er(err)
	}

	return pr
}

func (r *ReleaseCommand) updateReleasePullRequestBodyWithTickets(prNum int) {
	_, _, err := r.client.PullRequests.Edit(
		r.ctx,
		r.owner,
		r.repository,
		prNum,
		&github.PullRequest{
			Body: r.newPullRequestBodyFromCommits(prNum),
		})

	if err != nil {
		er(err)
	}
}

func (r *ReleaseCommand) newPullRequestBodyFromCommits(prNum int) *string {
	cmts, _, err := r.client.PullRequests.ListCommits(r.ctx, r.owner, r.repository, prNum, nil)
	if err != nil {
		er(err)
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

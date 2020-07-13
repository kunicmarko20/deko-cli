package cmd

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/kunicmarko20/deko-cli/util"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"time"
)

var (
	releaseCmd = &cobra.Command{
		Use:     "release",
		Aliases: []string{"r"},
		Short:   "Creates release branch and a PR in current git repository",
		Run:     NewReleaseCommand().Run,
	}
)

type ReleaseCommand struct {
	c  *util.GithubClient
	pb util.ProgressBar
}

func NewReleaseCommand() *ReleaseCommand {
	return &ReleaseCommand{
		pb: util.NewProgressBar(),
	}
}

func (r *ReleaseCommand) Run(cmd *cobra.Command, args []string) {
	//Needs to be initialized here because config gets loaded at this point
	r.c = util.NewGithubClient()

	if !prompter.YN(fmt.Sprintf("You are about to do a release for '%s' repository, are you sure?", r.c.Repository), true) {
		return
	}

	rlsDate := time.Now().Format("20060102")
	rlsBranch := "release-" + rlsDate

	r.pb.Render()

	r.createReleaseBranch(rlsBranch)
	pr := r.c.CreateReleasePullRequest(r.pb, rlsDate, rlsBranch)
	r.c.UpdateReleasePullRequestBodyWithTickets(r.pb, *pr.Number)

	r.pb.StopRendering()
	r.openPullRequestInBrowser(*pr.HTMLURL)
}

func (r *ReleaseCommand) createReleaseBranch(rlsBranch string) {
	util.CheckoutBranch(r.pb, "master")
	util.PullNewChanges(r.pb)
	util.CreateNewBranch(r.pb, rlsBranch)
	util.MergeChangesFromBranch(r.pb, "origin/staging")
	util.PushNewBranchToRemote(r.pb, "origin", rlsBranch)
}

func (r *ReleaseCommand) openPullRequestInBrowser(prUrl string) {
	fmt.Println(fmt.Sprintf("done, visit %s to see your pull request", prUrl))

	err := webbrowser.Open(prUrl)

	if err != nil {
		util.Exit(err)
	}
}

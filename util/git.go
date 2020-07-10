package util

import (
	"fmt"
	"github.com/ldez/go-git-cmd-wrapper/checkout"
	"github.com/ldez/go-git-cmd-wrapper/config"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/merge"
	"github.com/ldez/go-git-cmd-wrapper/push"
)

func ConfigGet(name string) string {
	url, err := git.Config(config.Get(name, ""))

	if err != nil {
		Exit(err)
	}

	return url
}
func CheckoutBranch(pb ProgressBar, branch string) {
	t := pb.Start("changing branch to master")

	msg, err := git.Checkout(checkout.Branch(branch))
	if err != nil {
		Exit(msg)
	}

	t.Increment(1)
}

func PullNewChanges(pb ProgressBar) {
	t := pb.Start("pulling new changes")
	msg, err := git.Pull()
	if err != nil {
		Exit(msg)
	}

	t.Increment(1)
}

func CreateNewBranch(pb ProgressBar, branch string) {
	t := pb.Start(fmt.Sprintf("creating new branch '%s'", branch))

	msg, err := git.Checkout(checkout.NewBranch(branch))
	if err != nil {
		Exit(msg)
	}

	t.Increment(1)
}

func MergeChangesFromBranch(pb ProgressBar, branch string) {
	t := pb.Start("merging " + branch)

	msg, err := git.Merge(merge.Commits(branch))
	if err != nil {
		Exit(msg)
	}

	t.Increment(1)
}

func PushNewBranchToRemote(pb ProgressBar, remote string, branch string) {
	t := pb.Start(fmt.Sprintf("pushing new branch '%s' to origin", branch))

	msg, err := git.Push(push.Remote(remote), push.Remote(branch))
	if err != nil {
		Exit(msg)
	}

	t.Increment(1)
}

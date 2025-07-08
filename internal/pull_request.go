package internal

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v31/github"
	"github.com/posener/goaction"
	"github.com/sirupsen/logrus"
	"strings"
)

type PullRequestCommitFetcher struct {
	pullRequest *github.PullRequestEvent
}

func (f *PullRequestCommitFetcher) getPullRequest() error {
	p, err := goaction.GetPullRequest()
	if err == nil {
		f.pullRequest = p
	}
	return err
}

func (f *PullRequestCommitFetcher) GetInfo() CommitsInfo {
	_ = f.getPullRequest()
	return CommitsInfo{
		Type:        "Pull Request",
		Target:      f.pullRequest.GetPullRequest().GetURL(),
		AuthorEmail: f.pullRequest.GetPullRequest().GetUser().GetEmail(),
		AuthorName:  f.pullRequest.GetPullRequest().GetUser().GetName(),
	}
}

func (f *PullRequestCommitFetcher) Test() bool {
	return f.getPullRequest() == nil
}

func (*PullRequestCommitFetcher) GetCommits() ([]Commit, error) {
	var commitMessages []Commit

	pr, _ := goaction.GetPullRequest()

	logrus.Debug("Opening local repository")
	repository, err := git.PlainOpen(".")
	if err != nil {
		return nil, fmt.Errorf("error opening local repository: %v", err)
	}

	logrus.Debug("Opening worktree")
	worktree, err := repository.Worktree()
	if err != nil {
		return nil, fmt.Errorf("error opening worktree: %v", err)
	}

	var branchName string
	if b, found := strings.CutPrefix(pr.GetPullRequest().GetHead().GetRef(), "refs/heads/"); found {
		branchName = b
	} else {
		return nil, fmt.Errorf("head reference %s doesn't include refs/heads/ prefix", pr.GetPullRequest().GetHead().GetRef())
	}

	logrus.Debugf("Checking out branch %s", branchName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		return nil, fmt.Errorf("error checking out branch %s: %v", branchName, err)
	}

	logrus.Debugf("Fetching git log for branch")
	log, err := repository.Log(&git.LogOptions{
		From: plumbing.NewHash(pr.GetPullRequest().GetHead().GetSHA()),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching git log: %v", err)
	}

	var skip bool
	err = log.ForEach(func(commit *object.Commit) error {
		if commit.Hash == plumbing.NewHash(pr.GetPullRequest().GetBase().GetSHA()) {
			skip = true
		}
		if skip {
			return nil
		}
		commitMessages = append(commitMessages, Commit{
			Message:     commit.Message,
			AuthorName:  commit.Author.Name,
			AuthorEmail: commit.Author.Email,
		})
		return nil
	})

	return commitMessages, nil
}

var _ CommitsFetcher = &PullRequestCommitFetcher{}

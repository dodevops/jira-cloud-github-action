package internal

import (
	"github.com/google/go-github/v31/github"
	"github.com/posener/goaction"
)

type PushCommitFetcher struct {
	push *github.PushEvent
}

func (f *PushCommitFetcher) getPush() error {
	p, err := goaction.GetPush()
	if err == nil {
		f.push = p
	}
	return err
}

func (f *PushCommitFetcher) GetInfo() CommitsInfo {
	_ = f.getPush()
	return CommitsInfo{
		Type:        "Push",
		Target:      f.push.GetRef(),
		AuthorEmail: f.push.GetPusher().GetEmail(),
		AuthorName:  f.push.GetPusher().GetName(),
	}
}

func (f *PushCommitFetcher) Test() bool {
	return f.getPush() == nil
}

func (*PushCommitFetcher) GetCommits() ([]Commit, error) {
	var commitMessages []Commit
	push, _ := goaction.GetPush()
	for _, commit := range push.Commits {
		commitMessages = append(commitMessages, Commit{
			Message:     commit.GetMessage(),
			AuthorName:  commit.Author.GetName(),
			AuthorEmail: commit.Author.GetEmail(),
		})
	}
	return commitMessages, nil
}

var _ CommitsFetcher = &PushCommitFetcher{}

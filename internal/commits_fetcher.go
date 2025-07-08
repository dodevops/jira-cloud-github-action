package internal

type CommitsInfo struct {
	Type        string
	Target      string
	AuthorEmail string
	AuthorName  string
}

type CommitsFetcher interface {
	GetInfo() CommitsInfo
	Test() bool
	GetCommits() ([]Commit, error)
}

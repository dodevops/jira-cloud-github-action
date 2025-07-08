package internal

type Commit struct {
	Message     string
	AuthorName  string
	AuthorEmail string
}

type CommentInfo struct {
	CommitsInfo CommitsInfo
	Commits     []Commit
}

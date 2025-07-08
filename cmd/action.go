package main

import (
	"bytes"
	"jira-cloud-github-action/internal"
	"log"
	"regexp"
	"slices"
	"text/template"

	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

func main() {
	var args struct {
		LogLevel          string   `arg:"env:INPUT_LOGLEVEL" help:"The loglevel to use" default:"INFO"`
		URL               string   `arg:"required,env:INPUT_URL" help:"The URL to your Jira Cloud instance"`
		User              string   `arg:"required,env:INPUT_USERNAME" help:"The user name for API requests"`
		Token             string   `arg:"required,env:INPUT_TOKEN" help:"The user's token for API requests"`
		CommitParseRegExp string   `arg:"env:INPUT_COMMITPARSEREGEXP" help:"The Regexp string to use to parse Jira issue ids from commit messages. Requires a named group called id that catches the whole issue id and a group project that catches only the project key" default:"\\((?P<id>(?P<project>[A-Za-z]+)-[0-9]+)\\)"`
		CommentTemplate   string   `arg:"env:INPUT_COMMENTTEMPLATE" help:"The go template for the issue comment. Takes the CommentInfo struct as an input" default:"{{ .CommitsInfo.Type }}"`
		DryRun            bool     `arg:"env:INPUT_DRYRUN" help:"Don't actually do something, just write what would be done'" default:"false"`
		OnlyProjects      []string `arg:"env:INPUT_ONLYPROJECTS" help:"Only allow commenting on these projects"`
	}
	arg.MustParse(&args)
	if l, err := logrus.ParseLevel(args.LogLevel); err != nil {
		log.Fatal(err)
	} else {
		logrus.SetLevel(l)
	}

	var matcher *regexp.Regexp
	if m, err := regexp.Compile(args.CommitParseRegExp); err != nil {
		logrus.Fatalf("Can not parse regular expression %s: %v", args.CommitParseRegExp, err)
	} else {
		matcher = m
	}

	api := internal.NewJiraAPI(args.URL, args.User, args.Token)

	var commentTemplate *template.Template
	if t, err := template.New("comment").Parse(args.CommentTemplate); err != nil {
		log.Fatalf("Can not parse template %s: %v", args.CommentTemplate, err)
	} else {
		commentTemplate = t
	}

	commitsFetchers := []internal.CommitsFetcher{&internal.PushCommitFetcher{}, &internal.PullRequestCommitFetcher{}}

	for _, fetcher := range commitsFetchers {
		if fetcher.Test() {
			logrus.Infof("Found %s", fetcher.GetInfo().Type)
			if commits, err := fetcher.GetCommits(); err != nil {
				log.Fatal(err)
			} else {
				for _, commit := range commits {
					matches := matcher.FindStringSubmatch(commit.Message)
					if len(matches) > 0 {
						project := matches[matcher.SubexpIndex("project")]
						if len(args.OnlyProjects) > 0 && !slices.Contains(args.OnlyProjects, project) {
							continue
						}
						issueID := matches[matcher.SubexpIndex("id")]
						commentInfo := internal.CommentInfo{
							CommitsInfo: fetcher.GetInfo(),
							Commits:     commits,
						}
						var commentWriter = bytes.Buffer{}
						if err := commentTemplate.Execute(&commentWriter, commentInfo); err != nil {
							log.Fatalf("Can not execute template %s with object %v: %v", args.CommentTemplate, commentInfo, err)
						}
						if err := api.CommentWorkItem(issueID, commentWriter.String()); err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}
}

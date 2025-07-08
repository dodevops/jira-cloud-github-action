package internal

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/summonio/markdown-to-adf/renderer"
)

type JiraAPI struct {
	client *resty.Client
}

func NewJiraAPI(url string, username string, token string) JiraAPI {
	return JiraAPI{
		client: resty.New().SetBaseURL(url).SetBasicAuth(username, token),
	}
}

func (a JiraAPI) CommentWorkItem(workItemID string, comment string) error {
	var adfComment bytes.Buffer
	if err := renderer.Render(&adfComment, []byte(comment)); err != nil {
		return fmt.Errorf("error compiling markdown to adf for comment %s: %v", comment, err)
	}
	logrus.Infof("Adding comment to issue %s: %s", workItemID, comment)
	if res, err := a.client.R().
		SetPathParam("issueIdOrKey", workItemID).
		SetBody(fmt.Sprintf("{ \"body\": %s }", adfComment.String())).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Post("/rest/api/3/issue/{issueIdOrKey}/comment"); err != nil {
		return fmt.Errorf("error adding comment: %v", err)
	} else {
		if res.IsError() {
			return fmt.Errorf("error adding comment (%d): %s", res.StatusCode(), res.String())
		}
	}
	return nil
}

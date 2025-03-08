package util

import (
	"context"
	"os"

	"github.com/google/go-github/v69/github"
)

var ghClient *github.Client

var owner string
var repo string

func init() {
	ghClient = github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	owner = os.Getenv("GITHUB_OWNER")
	repo = os.Getenv("GITHUB_REPO")
}

func CreateIssue(issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return ghClient.Issues.Create(context.Background(), owner, repo, issue)
}

func CreateIssueComment(number int, issueComment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return ghClient.Issues.CreateComment(context.Background(), owner, repo, number, issueComment)
}

func CloseIssue(number int) (*github.Issue, *github.Response, error) {
	state := "closed"
	return ghClient.Issues.Edit(context.Background(), owner, repo, number, &github.IssueRequest{
		State: &state,
	})
}

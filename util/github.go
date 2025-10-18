package util

import (
	"context"
	"log"

	"github.com/google/go-github/v69/github"
	"github.com/spf13/viper"
)

var ghClient *github.Client

var owner string
var repo string

func InitGithubClient() {
	ghClient = github.NewClient(nil).WithAuthToken(viper.GetString("github.token"))
	owner = viper.GetString("github.owner")
	repo = viper.GetString("github.repo")
	log.Println("GitHub client initialized for repo:", owner+"/"+repo)
}

func CreateIssue(issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return ghClient.Issues.Create(context.Background(), owner, repo, issue)
}

func CreateIssueComment(number int, issueComment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return ghClient.Issues.CreateComment(context.Background(), owner, repo, number, issueComment)
}

func GetIssue(number int) (*github.Issue, *github.Response, error) {
	return ghClient.Issues.Get(context.Background(), owner, repo, number)
}

func EditIssue(number int, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return ghClient.Issues.Edit(context.Background(), owner, repo, number, issue)
}

func ReactToIssueComment(commentId int64, reaction string) (*github.Reaction, *github.Response, error) {
	return ghClient.Reactions.CreateIssueCommentReaction(context.Background(), owner, repo, commentId, reaction)
}

func EditIssueComment(issueComment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return ghClient.Issues.EditComment(context.Background(), owner, repo, issueComment.GetID(), issueComment)
}

func CloseIssue(number int, stateReason string) (*github.Issue, *github.Response, error) {
	state := "closed"
	return ghClient.Issues.Edit(context.Background(), owner, repo, number, &github.IssueRequest{
		State:       &state,
		StateReason: &stateReason,
	})
}

func GetUserById(id int64) (*github.User, *github.Response, error) {
	return ghClient.Users.GetByID(context.Background(), id)
}

func AddIssueAssignee(issueNumber int, assignees []string) (*github.Issue, *github.Response, error) {
	return ghClient.Issues.AddAssignees(context.Background(), owner, repo, issueNumber, assignees)
}

func AddIssueLabels(issueNumber int, labels []string) ([]*github.Label, *github.Response, error) {
	return ghClient.Issues.AddLabelsToIssue(context.Background(), owner, repo, int(issueNumber), labels)
}

func RemoveIssueLabel(issueNumber int, label string) (*github.Response, error) {
	return ghClient.Issues.RemoveLabelForIssue(context.Background(), owner, repo, int(issueNumber), label)
}

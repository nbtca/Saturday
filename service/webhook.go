package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/webhooks/v6/github"
	githubapi "github.com/google/go-github/v69/github"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"
)

type GithubWebHook struct {
	hook *github.Webhook
}

func MakeGithubWebHook(secret string) (*GithubWebHook, error) {
	hook, err := github.New(github.Options.Secret(secret))
	if err != nil {
		return nil, err
	}
	return &GithubWebHook{hook: hook}, nil
}

func ExtractCommand(comment string) (string, error) {
	re := regexp.MustCompile(`@nbtca-bot\s+(\w+)`)
	match := re.FindStringSubmatch(comment)
	if len(match) < 2 {
		return "", fmt.Errorf("command not found")
	}
	return match[1], nil
}

func ExtractSizeLabel(label string) (string, error) {
	re := regexp.MustCompile(`size:\s*(\w+)`)
	match := re.FindStringSubmatch(label)
	if len(match) < 2 {
		return "", fmt.Errorf("size label not found")
	}
	err := ValidateEventSize(match[1])
	return match[1], err
}

func (gh *GithubWebHook) Handle(request *http.Request) error {
	payload, err := gh.hook.Parse(
		request,
		github.ReleaseEvent,
		github.PullRequestEvent,
		github.IssueCommentEvent,
		github.IssuesEvent,
	)
	if err != nil {
		return err
	}

	switch payload := payload.(type) {
	case github.IssuesPayload:
		return gh.handleIssuesPayload(payload)
	case github.IssueCommentPayload:
		return gh.handleIssueCommentPayload(payload)
	}
	return nil
}

func (gh *GithubWebHook) handleIssuesPayload(issue github.IssuesPayload) error {
	if issue.Action != "labeled" {
		return gh.handleIssueUnLabel(issue)
	}
	return gh.handleIssueWithLabel(issue)
}

func (gh *GithubWebHook) handleIssueUnLabel(issue github.IssuesPayload) error {
	size := ""
	for _, label := range issue.Issue.Labels {
		size, _ = ExtractSizeLabel(label.Name)
	}
	if size != "" {
		return nil
	}
	event, err := repo.GetEventByIssueId(issue.Issue.ID)
	if err != nil || event.EventId == 0 || event.Size == "" {
		return nil
	}
	event.Size = ""
	err = repo.UpdateEventSize(event.EventId, "")
	if err != nil {
		return fmt.Errorf("failed to update event size: %v", err)
	}
	_, _, err = RenderEventToGithubIssue(&event, int(issue.Issue.Number), &githubapi.IssueRequest{})
	util.Logger.Errorf("failed to render event to github issue: %v", err)
	return nil
}

func (gh *GithubWebHook) handleIssueWithLabel(issue github.IssuesPayload) error {
	if issue.Label == nil || issue.Label.ID == 0 {
		util.Logger.Debugf("issue label not found")
		return nil
	}
	size, err := ExtractSizeLabel(issue.Label.Name)
	if err != nil {
		util.Logger.Debugf(err.Error())
		return nil
	}
	util.Logger.Debugf("size label found: %s", size)
	event, err := repo.GetEventByIssueId(issue.Issue.ID)
	if err != nil {
		return fmt.Errorf("failed to get event by issue id: %v", err)
	}
	if err := repo.UpdateEventSize(event.EventId, size); err != nil {
		return fmt.Errorf("failed to update event size: %v", err)
	}
	event.Size = size
	_, _, err = RenderEventToGithubIssue(&event, int(issue.Issue.Number), &githubapi.IssueRequest{})
	util.Logger.Errorf("failed to render event to github issue: %v", err)
	return nil
}

func (gh *GithubWebHook) handleIssueCommentPayload(comment github.IssueCommentPayload) error {
	command, err := ExtractCommand(comment.Comment.Body)
	if err != nil {
		return err
	}
	util.Logger.Debugf("command from github webhook: %s", command)

	event, err := repo.GetEventByIssueId(comment.Issue.ID)
	if err != nil {
		return err
	}
	if event.EventId == 0 {
		return nil
	}
	member, err := MemberServiceApp.GetMemberByGithubId(strconv.FormatInt(comment.Sender.ID, 10))
	if err != nil {
		return err
	}
	if member.MemberId == "" {
		return util.MakeValidationError("member not found", nil)
	}
	util.Logger.Tracef("member found: %v", member)
	logtoUserRoleResponse, err := LogtoServiceApp.FetchUserRole(member.LogtoId)
	if err != nil {
		return fmt.Errorf("logto user role error %v", err)
	}
	identity := model.Identity{
		Id:     member.MemberId,
		Member: member,
		Role:   MemberServiceApp.MapLogtoUserRole(logtoUserRoleResponse),
	}
	util.Logger.Tracef("using identity %v", identity)

	return gh.processCommand(comment, command, event, identity)
}

func (gh *GithubWebHook) processCommand(comment github.IssueCommentPayload, command string, event model.Event, identity model.Identity) error {
	var readyForReviewLabel = "ready for review"
	var acceptedLabel = "accepted"

	var err error
	switch {
	case comment.Action == "created" && command == "accept":
		err = gh.handleAcceptCommand(comment, event, identity, acceptedLabel)
	case command == "commit":
		err = gh.handleCommitCommand(comment, event, identity, readyForReviewLabel)
	case comment.Action == "created" && command == "reject":
		err = gh.handleRejectCommand(comment, event, identity, readyForReviewLabel)
	case comment.Action == "created" && command == "close":
		err = EventServiceApp.Act(&event, identity, util.Close)
	case comment.Action == "created" && command == "drop":
		err = gh.handleDropCommand(comment, event, identity, acceptedLabel)
	}
	if err != nil {
		util.ReactToIssueComment(comment.Comment.ID, "confused")
	} else {
		util.ReactToIssueComment(comment.Comment.ID, "+1")
	}
	return err
}

func (gh *GithubWebHook) handleAcceptCommand(comment github.IssueCommentPayload, event model.Event, identity model.Identity, acceptedLabel string) error {
	err := EventServiceApp.Act(&event, identity, util.Accept)
	if err != nil {
		return err
	}
	user, _, err := util.GetUserById(comment.Sender.ID)
	if err != nil {
		return err
	}
	_, _, err = util.AddIssueAssignee(int(comment.Issue.Number), []string{*user.Login})
	if err != nil {
		return err
	}
	_, _, err = util.AddIssueLabels(int(comment.Issue.Number), []string{acceptedLabel})
	return err
}

func (gh *GithubWebHook) handleCommitCommand(comment github.IssueCommentPayload, event model.Event, identity model.Identity, readyForReviewLabel string) error {
	re := regexp.MustCompile(`@nbtca-bot\s+\w+`)
	text := comment.Comment.Body
	cleaned := re.ReplaceAllString(text, "")
	cleaned = strings.TrimSpace(cleaned)

	switch event.Status {
	case util.Accepted:
		if err := EventServiceApp.Act(&event, identity, util.Commit, cleaned); err != nil {
			return err
		}
		_, _, err := util.AddIssueLabels(int(comment.Issue.Number), []string{readyForReviewLabel})
		return err
	case util.Committed:
		return EventServiceApp.Act(&event, identity, util.AlterCommit, cleaned)
	}
	return nil
}

func (gh *GithubWebHook) handleRejectCommand(comment github.IssueCommentPayload, event model.Event, identity model.Identity, readyForReviewLabel string) error {
	err := EventServiceApp.Act(&event, identity, util.Reject)
	if err != nil {
		return err
	}
	// _, err = util.RemoveIssueLabel(int(comment.Issue.Number), readyForReviewLabel)
	return err
}

func (gh *GithubWebHook) handleDropCommand(comment github.IssueCommentPayload, event model.Event, identity model.Identity, acceptedLabel string) error {
	err := EventServiceApp.Act(&event, identity, util.Drop)
	if err != nil {
		return err
	}
	// _, err = util.RemoveIssueLabel(int(comment.Issue.Number), acceptedLabel)
	return err
}

type LogtoWebHook struct {
}

type UserEvent struct {
	Event        string         `json:"event"`
	CreatedAt    string         `json:"createdAt"`
	UserAgent    string         `json:"userAgent"`
	IP           string         `json:"ip"`
	Path         string         `json:"path"`
	Method       string         `json:"method"`
	Status       int            `json:"status"`
	Params       map[string]any `json:"params"`
	MatchedRoute string         `json:"matchedRoute"`
	Data         map[string]any `json:"data"`
	User         map[string]any `json:"user"`
	HookID       string         `json:"hookId"`
}

func (l *LogtoWebHook) Handle(request *http.Request) error {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}
	defer request.Body.Close() // Always close the body when done

	userEvent := &UserEvent{}
	if err := json.Unmarshal(bodyBytes, userEvent); err != nil {
		return fmt.Errorf("failed to parse request body into UserEvent: %v", err)
	}

	log.Printf("Received UserEvent: %+v", userEvent)

	var logtoUsersResponse FetchLogtoUsersResponse
	var user map[string]any
	if userEvent.Event == "PostSignIn" {
		user = userEvent.User
	} else {
		user = userEvent.Data
	}
	dataBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal UserEvent.Data: %v", err)
	}
	if err := json.Unmarshal(dataBytes, &logtoUsersResponse); err != nil {
		log.Printf("Not FetchLogtoUsersResponse: %v", err)
		return nil
	} else {
		log.Printf("Successfully mapped to FetchLogtoUsersResponse: %+v", logtoUsersResponse)
	}
	member, err := MemberServiceApp.GetMemberByLogtoId(logtoUsersResponse.Id)
	if err != nil {
		return err
	}
	if member.MemberId == "" {
		return util.MakeValidationError("member not found", nil)
	}
	member.Avatar = logtoUsersResponse.Avatar
	if logtoUsersResponse.Name != "" {
		member.Alias = logtoUsersResponse.Name
	}
	if logtoUsersResponse.UserName != "" {
		member.Alias = logtoUsersResponse.UserName
	}
	if gh, ok := logtoUsersResponse.Identities["github"]; ok {
		member.GithubId = gh.UserId
	}
	return MemberServiceApp.UpdateMember(member)

}

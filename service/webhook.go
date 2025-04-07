package service

import (
	"log"
	"net/http"

	"github.com/go-playground/webhooks/v6/github"
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

// - accept repair event when some one is assigned to the issue
// - close repair whe issue when rep
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
	switch payload.(type) {
	case github.IssuesPayload:
		issue := payload.(github.IssuesPayload)
		event, err := repo.GetEventByIssueId(issue.Issue.ID)
		if err != nil {
			return err
		}
		if event.EventId == 0 {
			return nil
		}
		log.Printf("event found %v", event)

		if issue.Action == "assigned" {
			return gh.onAssign(issue, event)
		}
		if issue.Action == "unassigned" {
			log.Printf("issue unassigned %v", issue)
		}
		if issue.Action == "closed" {
			log.Printf("issue closed %v", issue)
		}
	case github.IssueCommentPayload:
		comment := payload.(github.IssueCommentPayload)
		log.Printf("issue comment %+v", comment)
	}
	return nil
}

// assignee Id -> github user email -> logto user -> member
func (gh *GithubWebHook) onAssign(issue github.IssuesPayload, event model.Event) error {
	log.Printf("issue assigned %v", issue.Issue.ID)

	assigneeId := issue.Assignee.ID
	if assigneeId == 0 {
		return nil
	}
	assignee, _, _ := util.GetUserById(assigneeId)
	if assignee == nil {
		return nil
	}
	log.Printf("assignee %v", assignee)
	users, _ := LogtoServiceApp.FetchUsers(FetchLogtoUsersRequest{
		PageSize: 1,
		SearchParams: map[string]interface{}{
			"search.primaryEmail": *assignee.Email,
		},
	})
	if len(users) == 0 {
		return nil
	}
	user := users[0]
	member, err := MemberServiceApp.GetMemberByLogtoId(user.Id)
	if err != nil {
		return err
	}
	if member.MemberId == "" {
		return nil
	}

	token, _ := LogtoServiceApp.getToken()
	roles, _ := LogtoServiceApp.FetchUserRole(user.Id, token)
	log.Printf("member %v", member)
	err = EventServiceApp.Accept(&event, model.Identity{
		Id:     member.MemberId,
		Member: member,
		Role:   MemberServiceApp.MapLogtoUserRole(roles),
	})
	if err != nil {
		return err
	}
	return nil
}

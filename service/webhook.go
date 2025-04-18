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
	case github.IssueCommentPayload:
		comment := payload.(github.IssueCommentPayload)

		match := regexp.MustCompile(`@nbtca-bot\s+(\w+)`).FindStringSubmatch(comment.Comment.Body)
		if len(match) < 2 {
			return nil
		}
		action := match[1]

		event, err := repo.GetEventByIssueId(comment.Issue.ID)
		if err != nil {
			return err
		}
		if event.EventId == 0 {
			return nil
		}
		log.Printf("event found %v", event)
		member, err := MemberServiceApp.GetMemberByGithubId(strconv.FormatInt(comment.Sender.ID, 10))
		if err != nil {
			return err
		}
		if member.MemberId == "" {
			return util.MakeValidationError("member not found", nil)
		}
		logtoToken, _ := LogtoServiceApp.getToken()
		logtoUserRoleResponse, err := LogtoServiceApp.FetchUserRole(member.LogtoId, logtoToken)
		if err != nil {
			return fmt.Errorf("logto user role error %v", err)
		}
		identity := model.Identity{
			Id:     member.MemberId,
			Member: member,
			Role:   MemberServiceApp.MapLogtoUserRole(logtoUserRoleResponse),
		}

		if comment.Action == "created" && action == "accept" {
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

			_, _, err = util.AddIssueLabels(int(comment.Issue.Number), []string{"accepted"})
			if err != nil {
				return err
			}
		}
		if comment.Action == "created" && action == "commit" {
			re := regexp.MustCompile(`@nbtca-bot\s+\w+`)
			text := comment.Comment.Body
			cleaned := re.ReplaceAllString(text, "")
			cleaned = strings.TrimSpace(cleaned)

			if err := EventServiceApp.Act(&event, identity, util.Commit, cleaned); err != nil {
				return err
			}

			_, _, err := util.AddIssueLabels(int(comment.Issue.Number), []string{"ready for review"})
			if err != nil {
				return err
			}
		}

		if comment.Action == "created" && action == "close" {
			return EventServiceApp.Act(&event, identity, util.Close)
		}
	}
	return nil
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
	HookID       string         `json:"hookId"`
}

func (l *LogtoWebHook) Handle(request *http.Request) error {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	defer request.Body.Close() // Always close the body when done

	userEvent := &UserEvent{}
	if err := json.Unmarshal(bodyBytes, userEvent); err != nil {
		return fmt.Errorf("failed to parse request body into UserEvent: %v", err)
	}

	log.Printf("Received UserEvent: %+v", userEvent)

	// Try to map UserEvent.Data to FetchLogtoUsersResponse
	var logtoUsersResponse FetchLogtoUsersResponse
	dataBytes, err := json.Marshal(userEvent.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal UserEvent.Data: %v", err)
	}
	if err := json.Unmarshal(dataBytes, &logtoUsersResponse); err != nil {
		log.Printf("UserEvent.Data is not FetchLogtoUsersResponse: %v", err)
		return nil
	} else {
		log.Printf("UserEvent.Data successfully mapped to FetchLogtoUsersResponse: %+v", logtoUsersResponse)
	}
	member, err := MemberServiceApp.GetMemberByLogtoId(logtoUsersResponse.Id)
	if err != nil {
		return err
	}
	if member.MemberId == "" {
		return util.MakeValidationError("member not found", nil)
	}
	if gh, ok := logtoUsersResponse.Identities["github"]; ok {
		member.GithubId = gh.UserId
	} else {
		return nil
	}
	return MemberServiceApp.UpdateMember(member)

}

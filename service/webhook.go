package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nbtca/saturday/util"
)

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
		return err
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
	member.Alias = logtoUsersResponse.UserName
	if gh, ok := logtoUsersResponse.Identities["github"]; ok {
		member.GithubId = gh.UserId
	}
	return MemberServiceApp.UpdateMember(member)

}

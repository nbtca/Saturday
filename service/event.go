package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v69/github"
	md "github.com/nao1215/markdown"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"

	"gopkg.in/gomail.v2"
)

type EventService struct{}

func (service EventService) GetEventById(id int64) (model.Event, error) {
	event, err := repo.GetEventById(id)
	if err != nil {
		return model.Event{}, util.MakeInternalServerError()
	}
	if event.EventId == 0 {
		return model.Event{}, util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed")
	}
	return event, nil
}

func (service EventService) GetMemberEvents(f repo.EventFilter, memberId string) ([]model.Event, error) {
	return repo.GetMemberEvents(f, memberId)
}

func (service EventService) GetClientEvents(f repo.EventFilter, clientId string) ([]model.Event, error) {
	return repo.GetClientEvents(f, clientId)
}

func (service EventService) GetPublicEventById(id int64) (model.PublicEvent, error) {
	event, err := service.GetEventById(id)
	if err != nil {
		return model.PublicEvent{}, err
	}
	return model.CreatePublicEvent(event), nil
}

func (service EventService) GetPublicEvents(f repo.EventFilter) ([]model.PublicEvent, error) {
	events, err := repo.GetEvents(f)
	if err != nil {
		return nil, err
	}
	publicEvents := make([]model.PublicEvent, len(events))
	for i, v := range events {
		publicEvents[i] = model.CreatePublicEvent(v)
	}
	return publicEvents, nil
}

func (service EventService) CreateEvent(event *model.Event) error {
	// insert event
	if err := repo.CreateEvent(event); err != nil {
		return err
	}
	identity := model.Identity{
		Id:   fmt.Sprint(event.ClientId),
		Role: "client",
	}
	// insert event status and event log
	if err := service.Act(event, identity, util.Create); err != nil {
		return err
	}
	return nil
}

func (service EventService) Accept(event *model.Event, identity model.Identity) error {
	if err := service.Act(event, identity, util.Accept); err != nil {
		return err
	}
	return nil
}

func (service EventService) SendActionNotify(event *model.Event, subject string) error {
	if event == nil {
		return util.MakeInternalServerError()
	}
	service.SendActionNotifyViaMail(event, subject)
	return nil
}

func (service EventService) SendActionNotifyViaMail(event *model.Event, subject string) error {
	m := gomail.NewMessage()
	receiverAddress := os.Getenv("MAIL_RECEIVER_ADDRESS")
	if receiverAddress == "" {
		return fmt.Errorf("MAIL_RECEIVER_ADDRESS is not set")
	}
	m.SetHeader("To", receiverAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", fmt.Sprintf(
		`<div>
  <span style="padding-right:10px;">型号:</span>
  <span>%s</span>
</div>
<div>
  <span style="padding-right:10px;">问题描述:</span>
  <span>%s</span>
</div>
<div>
  <span style="padding-right:10px;">创建时间:</span>
  <span>%s</span>
</div>
<div style="padding-top:10px;">
  <a href="https://repair.nbtca.space">在 Sunday 中处理</a>
</div>`, event.Model, event.Problem, event.GmtCreate))

	if err := util.SendMail(m); err != nil {
		return util.MakeInternalServerError().SetMessage("fail on mail")
	}
	return nil
}

func syncEventActionToGithubIssue(event *model.Event, eventLog model.EventLog, identity model.Identity) error {
	if util.Action(eventLog.Action) == util.Create {
		body := event.ToMarkdownString()
		title := fmt.Sprintf("%s(#%v)", event.Problem, event.EventId)
		issue, _, err := util.CreateIssue(&github.IssueRequest{
			Title:  &title,
			Body:   &body,
			Labels: &[]string{"ticket"},
		})
		if err != nil {
			return err
		}
		event.GithubIssueId = sql.NullInt64{
			Valid: true,
			Int64: int64(*issue.ID),
		}
		event.GithubIssueNumber = sql.NullInt64{
			Valid: true,
			Int64: int64(*issue.Number),
		}
		return nil
	}
	if !event.GithubIssueId.Valid {
		return fmt.Errorf("event.GithubIssueId is not valid")
	}

	buf := new(bytes.Buffer)
	description := md.NewMarkdown(buf).
		H2(eventLog.Action).
		PlainText(eventLog.Description)
	if util.Action(eventLog.Action) == util.Cancel {
		description = description.PlainText("Cancelled by client")
	} else {
		description = description.PlainText(fmt.Sprintf("By %s", identity.Member.Alias))
	}
	commentBody := description.String()

	_, _, err := util.CreateIssueComment(int(event.GithubIssueNumber.Int64), &github.IssueComment{
		Body: &commentBody,
	})
	if err != nil {
		return err
	}

	if util.Action(eventLog.Action) == util.Close {
		if _, _, err := util.CloseIssue(int(event.GithubIssueNumber.Int64), "completed"); err != nil {
			return err
		}
	} else if util.Action(eventLog.Action) == util.Cancel {
		if _, _, err := util.CloseIssue(int(event.GithubIssueNumber.Int64), "not_planned"); err != nil {
			return err
		}
	}
	return nil
}

/*
this function validates the action and then perform action to the event.
it also persists the event and event log.
*/
func (service EventService) Act(event *model.Event, identity model.Identity, action util.Action, description ...string) error {
	handler := util.MakeEventActionHandler(action, event, identity)
	if err := handler.ValidateAction(); err != nil {
		util.Logger.Error("validate action failed: ", err)
		return err
	}
	for _, d := range description {
		handler.Description = fmt.Sprint(handler.Description, d)
	}

	log := handler.Handle()

	err := syncEventActionToGithubIssue(event, log, identity)
	if err != nil {
		util.Logger.Error(err)
	}

	// persist event
	if err := repo.UpdateEvent(event, &log); err != nil {
		return err
	}
	// append log
	event.Logs = append(event.Logs, log)
	return nil
}

var EventServiceApp = EventService{}

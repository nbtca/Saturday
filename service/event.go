package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v69/github"
	md "github.com/nao1215/markdown"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"

	"gopkg.in/gomail.v2"
)

type EventService struct{}

func (service EventService) GetEventById(id int64) (model.Event, error) {
	event, err := repo.GetEventById(id)
	if err != nil {
		return model.Event{}, util.
			MakeServiceError(http.StatusInternalServerError).
			SetMessage(err.Error())
	}
	if event.EventId == 0 {
		return model.Event{}, util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed")
	}
	return event, nil
}

func (service EventService) ExportEventToXlsx(f repo.EventFilter, startTime, endTime string) (*excelize.File, error) {
	events, err := repo.GetClosedEventsByTimeRange(f, startTime, endTime)
	if err != nil {
		return nil, err
	}
	eventsExported := make([]model.EventExported, len(events))
	for i, v := range events {
		eventsExported[i] = model.EventExported{
			EventId:          v.Event.EventId,
			MemberId:         v.Member.MemberId.String,
			EventDescription: v.Event.Problem,
			MemberName:       v.Member.Name.String,
			MemberSection:    v.Member.Section.String,
			MemberPhone:      v.Member.Phone.String,
			EventSize:        v.Event.Size,
			EventStatus:      v.Event.Status,
			CreatedAt:        v.Event.GmtCreate,
			ClosedAt:         v.Event.GmtModified,
		}
		log.Println("event", v.Event)
		log.Println("closed_by", v.Admin)
		if v.Admin.Member().MemberId != "" {
			eventsExported[i].ClosedByMemberId = v.Admin.MemberId.String
		}
	}
	const MaxHour = 8
	groupedByMember := make(map[string]model.EventExportedGroupedByMember)
	for _, event := range eventsExported {
		memberId := event.MemberId
		if _, exists := groupedByMember[memberId]; !exists {
			groupedByMember[memberId] = model.EventExportedGroupedByMember{
				MemberId:      event.MemberId,
				MemberName:    event.MemberName,
				MemberSection: event.MemberSection,
				MemberPhone:   event.MemberPhone,
				Hour:          2, // Base hour for each member
			}
		}
		group := groupedByMember[memberId]
		sizeHour := EventSizeToHour(event.EventSize) // Increment hour count for each event
		if sizeHour > 0 {
			group.Hour += sizeHour
		} else {
			group.Hour += 0.5 // Default increment for unknown sizes
		}
		if group.Hour > MaxHour {
			group.Hour = MaxHour // Cap the hour count at MaxHour
		}
		groupedByMember[memberId] = group
	}

	result := make([]model.EventExportedGroupedByMember, 0, len(groupedByMember))
	for _, grouped := range groupedByMember {
		result = append(result, grouped)
	}

	excelFile := excelize.NewFile()
	defer func() {
		if err := excelFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	groupedByMemberSheet := "Sheet1"
	index, err := excelFile.NewSheet(groupedByMemberSheet)
	if err != nil {
		util.Logger.Error(err)
		return nil, err
	}
	excelFile.SetCellValue(groupedByMemberSheet, "A1", "学号")
	excelFile.SetCellValue(groupedByMemberSheet, "B1", "姓名")
	excelFile.SetCellValue(groupedByMemberSheet, "C1", "班级")
	excelFile.SetCellValue(groupedByMemberSheet, "D1", "联系方式")
	excelFile.SetCellValue(groupedByMemberSheet, "E1", "时长")
	for i, event := range result {
		excelFile.SetCellValue(groupedByMemberSheet, fmt.Sprintf("A%v", i+2), event.MemberId)
		excelFile.SetCellValue(groupedByMemberSheet, fmt.Sprintf("B%v", i+2), event.MemberName)
		excelFile.SetCellValue(groupedByMemberSheet, fmt.Sprintf("C%v", i+2), event.MemberSection)
		excelFile.SetCellValue(groupedByMemberSheet, fmt.Sprintf("D%v", i+2), event.MemberPhone)
		excelFile.SetCellValue(groupedByMemberSheet, fmt.Sprintf("E%v", i+2), event.Hour)
	}
	overAllSheet := "Sheet2"
	_, err = excelFile.NewSheet(overAllSheet)
	if err != nil {
		util.Logger.Error(err)
		return nil, err
	}
	excelFile.SetCellValue(overAllSheet, "A1", "学号")
	excelFile.SetCellValue(overAllSheet, "B1", "姓名")
	excelFile.SetCellValue(overAllSheet, "C1", "班级")
	excelFile.SetCellValue(overAllSheet, "D1", "事件编号")
	excelFile.SetCellValue(overAllSheet, "E1", "事件描述")
	excelFile.SetCellValue(overAllSheet, "F1", "工作量")
	excelFile.SetCellValue(overAllSheet, "G1", "事件状态")
	excelFile.SetCellValue(overAllSheet, "H1", "创建时间")
	excelFile.SetCellValue(overAllSheet, "I1", "关闭时间")
	excelFile.SetCellValue(overAllSheet, "J1", "审核人")
	excelFile.SetCellValue(overAllSheet, "K1", "GithubIssue")
	githubIssueBaseUrl := fmt.Sprintf("https://github.com/%v/%v/issues", viper.GetString("github.owner"), viper.GetString("GITHUB_REPO"))
	for i, event := range eventsExported {
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("A%v", i+2), event.MemberId)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("B%v", i+2), event.MemberName)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("C%v", i+2), event.MemberSection)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("D%v", i+2), event.EventId)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("E%v", i+2), event.EventDescription)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("F%v", i+2), event.EventSize)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("G%v", i+2), event.EventStatus)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("H%v", i+2), event.CreatedAt)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("I%v", i+2), event.ClosedAt)
		excelFile.SetCellValue(overAllSheet, fmt.Sprintf("J%v", i+2), event.ClosedByMemberId)
		if event.EventGithubIssueNumber != 0 {
			excelFile.SetCellValue(overAllSheet, fmt.Sprintf("K%v", i+2), fmt.Sprintf("%v/%v", githubIssueBaseUrl, event.EventGithubIssueNumber))
		}
	}

	excelFile.SetActiveSheet(index)
	return excelFile, nil
}

func (service EventService) GetMemberEvents(f repo.EventFilter, memberId string) ([]model.Event, error) {
	return repo.GetMemberEvents(f, memberId)
}

func (service EventService) GetClientEvents(f repo.EventFilter, clientId int64) ([]model.Event, error) {
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

func (service EventService) GetPublicEventsWithCount(f repo.EventFilter) ([]model.PublicEvent, int64, error) {
	events, err := service.GetPublicEvents(f)
	if err != nil {
		return nil, 0, err
	}
	count, err := repo.CountEvents(f)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}

func (service EventService) GetMemberEventsWithCount(f repo.EventFilter, memberId string) ([]model.Event, int64, error) {
	events, err := repo.GetMemberEvents(f, memberId)
	if err != nil {
		return nil, 0, err
	}
	count, err := repo.CountMemberEvents(f, memberId)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}

func (service EventService) GetClientEventsWithCount(f repo.EventFilter, clientId int64) ([]model.Event, int64, error) {
	events, err := repo.GetClientEvents(f, clientId)
	if err != nil {
		return nil, 0, err
	}
	count, err := repo.CountClientEvents(f, clientId)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}

func (service EventService) CreateEvent(event *model.Event) error {
	// insert event
	if err := repo.CreateEvent(event); err != nil {
		return err
	}
	identity := model.Identity{
		Id:       fmt.Sprint(event.ClientId),
		ClientId: event.ClientId,
		Role:     "client",
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

func (service EventService) SendActionNotify(event *model.Event, eventLog model.EventLog, identity model.Identity) {
	if event == nil {
		return
	}

	go func() {
		err := service.SendActionNotifyViaMail(event, eventLog, identity)
		if err != nil {
			util.Logger.Error("send action notify via mail failed: ", err)
		}
	}()
	go func() {
		err := service.SendActionNotifyViaNSQ(event, eventLog, identity)
		if err != nil {
			util.Logger.Error("send action notify via nsq failed: ", err)
		}
	}()

}

func (service EventService) SendActionNotifyViaNSQ(event *model.Event, eventLog model.EventLog, identity model.Identity) error {
	producer := util.GetNSQProducer()
	if producer == nil {
		return nil
	}
	mapEventLog := map[string]interface{}{
		"event_id":    eventLog.EventId,
		"member_id":   eventLog.MemberId,
		"action":      eventLog.Action,
		"problem":     event.Problem,
		"model":       event.Model,
		"gmt_create":  eventLog.GmtCreate,
		"description": eventLog.Description,
	}
	if identity.Member.Alias != "" {
		mapEventLog["member_alias"] = identity.Member.Alias
	} else {
		mapEventLog["member_alias"] = ""
	}
	jsonMap, _ := json.Marshal(mapEventLog)
	return producer.PublishAsync(util.EventTopic, jsonMap, nil)
}

func (service EventService) SendActionNotifyViaMail(event *model.Event, eventLog model.EventLog, identity model.Identity) error {
	switch eventLog.Action {
	case string(util.Accept):
		m := gomail.NewMessage()
		issueNumber, err := event.GithubIssueNumber.Value()
		if err != nil {
			return fmt.Errorf("event.GithubIssueNumber is not valid: %v", err)
		}
		if identity.Member.LogtoId == "" {
			return fmt.Errorf("identity.Member.LogtoId is not set")
		}
		logtoUser, err := LogtoServiceApp.FetchUserById(identity.Member.LogtoId)
		if err != nil {
			return fmt.Errorf("fetch logto user failed: %v", err)
		}

		m.SetHeader("To", logtoUser.PrimaryEmail)
		m.SetHeader("Subject", fmt.Sprintf("维修状态更新(#%v)", event.EventId))
		m.SetBody("text/html", fmt.Sprintf(
			`
					<h3>新的状态为: %v</h3>
		<div>
  		<span style="padding-right:10px;">问题描述:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">型号:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">创建时间:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">手机:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">QQ:</span>
  		<span>%s</span>
		</div>
		<div style="padding-top:10px;">
  		<a href="http://github.com/nbtca/repair-tickets/issues/%v">在 nbtca/repair-tickets 中处理</a>
		</div>
			`, event.Status, event.Problem, event.Model, event.GmtCreate, event.Phone, event.QQ, issueNumber))

		if err := util.SendMail(m); err != nil {
			return util.MakeInternalServerError().SetMessage("fail on mail")
		}
		util.Logger.Trace("event accepted, send mail to ", logtoUser.PrimaryEmail)
		return nil
	}
	return nil
}

type EventAnalyzeResult struct {
	Suggestion string
	Tag        string
}

func (service EventService) Analyze(event *model.Event) (EventAnalyzeResult, error) {
	request := WorkflowRunRequest{
		Inputs: map[string]interface{}{
			"EventId": event.EventId,
		},
		ResponseMode: "blocking",
		User:         "saturday",
	}
	response, err := RunDifyWorkflow(request)
	if err != nil {
		return EventAnalyzeResult{}, err
	}
	if response.Data.Error != nil {
		return EventAnalyzeResult{}, fmt.Errorf("error: %v", response.Data.Error)
	}
	if response.Data.Outputs == nil {
		return EventAnalyzeResult{}, fmt.Errorf("no outputs")
	}
	result := EventAnalyzeResult{
		Suggestion: response.Data.Outputs["suggestion"].(string),
		Tag:        response.Data.Outputs["tag"].(string),
	}
	return result, nil
}

func RenderEventToMarkdownString(event *model.Event) string {
	body := event.ToMarkdown()
	log.Println("member", event.Member)
	log.Println("closedBy", event.ClosedByMember)
	memberAlias := ""
	if event.Member != nil {
		memberAlias = event.Member.Alias
	}
	closedByAlias := ""
	if event.ClosedByMember != nil {
		closedByAlias = event.ClosedByMember.Alias
	}
	body.HorizontalRule()
	body.Table(md.TableSet{
		Header: []string{"Field", "Value", "Description"},
		Rows: [][]string{
			{"Current Status", event.Status, ""},
			{"Size", event.Size, ""},
			{"Accepted By", memberAlias, ""},
			{"Closed By", closedByAlias, ""},
		},
	})
	body.HorizontalRule()

	mermaidDiagram := `flowchart LR
	A[Open] --> |Drop| B[Canceled]
	A --> |Accept| C[Accepted]
	C --> |Commit| D[Committed]
	D --> |AlterCommit| D
	D --> |Approve| E[Closed]
	D --> |Reject| C`

	buf := new(bytes.Buffer)
	m := md.NewMarkdown(buf)
	m.LF()
	m.LF()
	m.PlainText("You can update event status by commenting on this Issue:")
	m.BulletList(
		"`@nbtca-bot accept` will accept this ticket",
		"`@nbtca-bot drop` will drop your previous accept",
		"`@nbtca-bot commit` will submit this ticket for admin approval",
		"`@nbtca-bot reject` will send this ticket back to assignee",
		"`@nbtca-bot close` will close this ticket as completed",
	)
	m.CodeBlocks(md.SyntaxHighlightMermaid, mermaidDiagram)
	m.Blockquote("Get more detailed documentation at [docs.nbtca.space/repair/weekend](http://docs.nbtca.space/repair/weekend.html)")

	body.Details("nbtca-bot commands", m.String())
	return body.String()
}

func RenderEventToGithubIssue(event *model.Event, issueNumber int, issueRequest *github.IssueRequest) (*github.Issue, *github.Response, error) {
	bodyString := RenderEventToMarkdownString(event)
	issueRequest.Body = &bodyString
	title := fmt.Sprintf("%s(#%v)", event.Problem, event.EventId)
	issueRequest.Title = &title
	return util.EditIssue(issueNumber, issueRequest)
}

func CreateEventAnalyzeComment(event *model.Event, issue *github.Issue) error {
	analyzeResult, err := EventServiceApp.Analyze(event)
	if err != nil {
		return fmt.Errorf("analyze event failed: %v", err)
	}
	_, _, err = util.CreateIssueComment(int(issue.GetNumber()), &github.IssueComment{
		Body: &analyzeResult.Suggestion,
	})
	if err != nil {
		return fmt.Errorf("create issue comment for event analyze failed: %v", err)
	}
	return nil
}

func syncEventActionToGithubIssue(event *model.Event, eventLog model.EventLog, identity model.Identity) error {
	if util.Action(eventLog.Action) == util.Create {
		bodyString := RenderEventToMarkdownString(event)
		title := fmt.Sprintf("%s(#%v)", event.Problem, event.EventId)
		issue, _, err := util.CreateIssue(&github.IssueRequest{
			Title:  &title,
			Body:   &bodyString,
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

		go func(event *model.Event, issue *github.Issue) {
			err := CreateEventAnalyzeComment(event, issue)
			if err != nil {
				util.Logger.Error("create event analyze comment failed: ", err)
			}
		}(event, issue)

		return nil
	}

	if !event.GithubIssueId.Valid {
		return fmt.Errorf("event.GithubIssueId is not valid")
	}

	RenderEventToGithubIssue(event, int(event.GithubIssueNumber.Int64), &github.IssueRequest{})

	buf := new(bytes.Buffer)
	memberName := identity.Member.Alias
	if identity.Member.LogtoId != "" {
		logtoUser, err := LogtoServiceApp.FetchUserById(identity.Member.LogtoId)
		if err != nil {
			util.Logger.Error("fetch logto user failed: ", err)
			return err
		}
		memberName = fmt.Sprintf("%v (%v)", logtoUser.Name, logtoUser.PrimaryEmail)
	}
	description := md.NewMarkdown(buf).
		H2(eventLog.Action).
		PlainText(eventLog.Description)
	if util.Action(eventLog.Action) == util.Cancel {
		description = description.PlainText("Cancelled by client")
	} else {
		description = description.PlainText(fmt.Sprintf("By %s", memberName))
	}
	commentBody := description.String()

	_, _, err := util.CreateIssueComment(int(event.GithubIssueNumber.Int64), &github.IssueComment{
		Body: &commentBody,
	})
	if err != nil {
		return err
	}

	var readyForReviewLabel = "ready for review"
	var acceptedLabel = "accepted"

	if util.Action(eventLog.Action) == util.Close {
		if _, _, err := util.CloseIssue(int(event.GithubIssueNumber.Int64), "completed"); err != nil {
			return err
		}
	} else if util.Action(eventLog.Action) == util.Cancel {
		if _, _, err := util.CloseIssue(int(event.GithubIssueNumber.Int64), "not_planned"); err != nil {
			return err
		}
	} else if util.Action(eventLog.Action) == util.Accept {
		_, _, err = util.AddIssueLabels(int(event.GithubIssueId.Int64), []string{acceptedLabel})
		if err != nil {
			util.Logger.Error("add issue labels failed: ", err)
		}
	} else if util.Action(eventLog.Action) == util.Commit {
		_, _, err = util.AddIssueLabels(int(event.GithubIssueId.Int64), []string{readyForReviewLabel})
		if err != nil {
			util.Logger.Error("add issue labels failed: ", err)
		}
	} else if util.Action(eventLog.Action) == util.Drop {
		_, err = util.RemoveIssueLabel(int(event.GithubIssueId.Int64), acceptedLabel)
		if err != nil {
			util.Logger.Error("remove issue labels failed: ", err)
		}
	} else if util.Action(eventLog.Action) == util.Reject {
		_, err = util.RemoveIssueLabel(int(event.GithubIssueId.Int64), readyForReviewLabel)
		if err != nil {
			util.Logger.Error("remove issue labels failed: ", err)
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

	// persist event
	if err := repo.UpdateEvent(event, &log); err != nil {
		return err
	}
	// append log
	event.Logs = append(event.Logs, log)

	err := syncEventActionToGithubIssue(event, log, identity)
	if err != nil {
		util.Logger.Error(err)
	}

	service.SendActionNotify(event, log, identity)
	util.Logger.Tracef("event log: %v", log)

	return nil
}

func ValidateEventSize(size string) error {
	// size can be one of xs,s,m,l,xl
	if size != "xs" && size != "s" && size != "m" && size != "l" && size != "xl" {
		return fmt.Errorf("size %s is not valid", size)
	}
	return nil
}

func EventSizeToHour(size string) float64 {
	// size can be one of xs,s,m,l,xl
	switch size {
	case "xs":
		return 0.5
	case "s":
		return 1
	case "m":
		return 2
	case "l":
		return 4
	case "xl":
		return 8
	default:
		return 0
	}
}

var EventServiceApp = EventService{}

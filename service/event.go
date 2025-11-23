package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

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
	var recipients []model.Member

	switch eventLog.Action {
	case string(util.Create):
		// Send to all members who have enabled new event notifications
		optedInMembers, err := MemberServiceApp.GetMembersWithNotificationEnabled(model.NotifNewEventCreated)
		if err != nil {
			util.Logger.Errorf("failed to get members with enabled notifications: %v", err)
			return nil
		}
		recipients = optedInMembers

	case string(util.Accept), string(util.Drop), string(util.Commit), string(util.AlterCommit):
		// Send to the acting member if they have event_assigned_to_me enabled
		if identity.Member.LogtoId != "" {
			member, err := MemberServiceApp.GetMemberByLogtoId(identity.Member.LogtoId)
			if err == nil {
				prefs := member.GetNotificationPreferences()
				if prefs.EventAssignedToMe {
					recipients = append(recipients, member)
				}
			}
		}

	case string(util.Cancel), string(util.Reject), string(util.Close):
		// Send to the assigned member if they have event_assigned_to_me enabled
		if event.MemberId != "" {
			member, err := MemberServiceApp.GetMemberById(event.MemberId)
			if err == nil {
				prefs := member.GetNotificationPreferences()
				if prefs.EventAssignedToMe {
					recipients = append(recipients, member)
				}
			}
		}
	}

	// Send emails to all recipients
	for _, recipient := range recipients {
		if recipient.LogtoId == "" {
			continue
		}

		// Fetch user email from Logto
		logtoUser, err := LogtoServiceApp.FetchUserById(recipient.LogtoId)
		if err != nil {
			util.Logger.Errorf("fetch logto user failed for %s: %v", recipient.LogtoId, err)
			continue
		}

		// Generate email subject and content based on action
		subject, bodyHTML := service.generateEmailContent(event, eventLog)

		// Create and send email
		m := gomail.NewMessage()
		m.SetHeader("To", logtoUser.PrimaryEmail)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", bodyHTML)

		if err := util.SendMail(m); err != nil {
			util.Logger.Errorf("failed to send email to %s: %v", logtoUser.PrimaryEmail, err)
			continue
		}

		util.Logger.Tracef("send email for action '%s' to %s", eventLog.Action, logtoUser.PrimaryEmail)
	}

	return nil
}

func getEventStatusText(status string) string {
	statusTextMap := map[string]string{
		util.Open:      "待处理",
		util.Accepted:  "维修中",
		util.Committed: "维修中",
		util.Closed:    "已完成",
		util.Cancelled: "已取消",
	}
	if text, ok := statusTextMap[status]; ok {
		return text
	}
	return status
}

func (service EventService) generateEmailContent(event *model.Event, eventLog model.EventLog) (string, string) {
	var actionTitle string
	var actionMessage string

	switch eventLog.Action {
	case string(util.Accept):
		actionTitle = "工单已接受"
		actionMessage = "您已接受此维修工单，请及时处理。"
	case string(util.Cancel):
		actionTitle = "工单已取消"
		actionMessage = "此维修工单已被客户取消。"
	case string(util.Drop):
		actionTitle = "工单已放弃"
		actionMessage = "您已放弃此维修工单，工单将重新开放。"
	case string(util.Commit):
		actionTitle = "工单已提交审核"
		actionMessage = "您已提交此维修工单等待审核。"
	case string(util.AlterCommit):
		actionTitle = "工单提交已更新"
		actionMessage = "您已更新此维修工单的提交内容。"
	case string(util.Reject):
		actionTitle = "工单已被驳回"
		actionMessage = "您提交的维修工单未通过审核，请修改后重新提交。"
		if eventLog.Description != "" {
			actionMessage = fmt.Sprintf("您提交的维修工单未通过审核，原因：%s", eventLog.Description)
		}
	case string(util.Close):
		actionTitle = "工单已完成"
		actionMessage = "恭喜！此维修工单已成功完成并关闭。"
	default:
		actionTitle = "工单状态更新"
		actionMessage = ""
	}

	subject := fmt.Sprintf("维修工单 #%v - %s", event.EventId, actionTitle)
	statusText := getEventStatusText(event.Status)

	// Build web URL with configurable hostname
	webHostname := viper.GetString("web.hostname")
	if webHostname == "" {
		webHostname = "nbtca.space"
	}

	// Build URL with status filter and event ID
	// Include all statuses in filter for better user experience
	statusFilter := url.QueryEscape("open,accepted,committed,closed")
	webURL := fmt.Sprintf("https://%s/repair/admin?page=1&status=%s&eventid=%d",
		webHostname, statusFilter, event.EventId)

	bodyHTML := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">%s</h2>
			<p style="color: #666;">%s</p>
			<div style="background-color: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0;">
				<h3 style="margin-top: 0; color: #333;">当前状态: %s</h3>
				<div style="margin: 10px 0;">
					<span style="font-weight: bold; color: #555;">问题描述:</span>
					<span style="color: #333;">%s</span>
				</div>
				<div style="margin: 10px 0;">
					<span style="font-weight: bold; color: #555;">型号:</span>
					<span style="color: #333;">%s</span>
				</div>
				<div style="margin: 10px 0;">
					<span style="font-weight: bold; color: #555;">创建时间:</span>
					<span style="color: #333;">%s</span>
				</div>
				<div style="margin: 10px 0;">
					<span style="font-weight: bold; color: #555;">联系方式:</span>
					<span style="color: #333;">手机: %s | QQ: %s</span>
				</div>
			</div>
			<div style="margin-top: 20px;">
				<a href="%s"
				   style="display: inline-block; padding: 10px 20px; background-color: #0366d6; color: white; text-decoration: none; border-radius: 5px;">
					查看工单详情
				</a>
			</div>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">
				这是一封自动发送的邮件，请勿直接回复。
			</p>
		</div>
	`, actionTitle, actionMessage, statusText, event.Problem, event.Model, event.GmtCreate, event.Phone, event.QQ, webURL)

	return subject, bodyHTML
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

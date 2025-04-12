package util

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nbtca/saturday/model"
	"github.com/sirupsen/logrus"
)

const (
	Open      string = "open"
	Cancelled string = "cancelled"
	Accepted  string = "accepted"
	Committed string = "committed"
	Closed    string = "closed"
)

type Action string

const (
	Create      Action = "create"
	Accept      Action = "accept"
	Cancel      Action = "cancel"
	Drop        Action = "drop"
	Commit      Action = "commit"
	AlterCommit Action = "alterCommit"
	Reject      Action = "reject"
	Close       Action = "close"
	Update      Action = "update"
)

type customLogFunc func(*eventActionHandler) model.EventLog

type eventActionHandler struct {
	event       *model.Event
	actor       model.Identity
	action      Action
	role        []string
	prevStatus  string
	nextStatus  string
	Description string
	customLog   customLogFunc
}

var idLog customLogFunc = func(eh *eventActionHandler) model.EventLog {
	return eh.createEventLog(createEventLogArgs{
		Id: eh.actor.Id,
	})
}

var idAndDescriptionLog customLogFunc = func(eh *eventActionHandler) model.EventLog {
	return eh.createEventLog(createEventLogArgs{
		Id:          eh.actor.Id,
		Description: eh.Description,
	})
}

var eventActionMap map[Action]eventActionHandler = map[Action]eventActionHandler{
	Create: {
		action:     Create,
		role:       []string{"client"},
		prevStatus: "",
		nextStatus: Open,
	},
	Accept: {
		action:     Accept,
		role:       []string{"member", "admin"},
		prevStatus: Open,
		nextStatus: Accepted,
		customLog: func(eh *eventActionHandler) model.EventLog {
			eh.event.MemberId = eh.actor.Id
			return eh.createEventLog(createEventLogArgs{
				Id: eh.actor.Id,
			})
		},
	},
	Cancel: {
		action:     Cancel,
		role:       []string{"client_current"},
		prevStatus: Open,
		nextStatus: Cancelled,
	},
	Drop: {
		action:     Drop,
		role:       []string{"member_current"},
		prevStatus: Accepted,
		nextStatus: Open,
		customLog: func(eh *eventActionHandler) model.EventLog {
			eh.event.MemberId = ""
			return eh.createEventLog(createEventLogArgs{
				Id: eh.actor.Id,
			})
		},
	},
	Commit: {
		action:     Commit,
		role:       []string{"member_current"},
		prevStatus: Accepted,
		nextStatus: Committed,
		customLog:  idAndDescriptionLog,
	},
	AlterCommit: {
		action:     AlterCommit,
		role:       []string{"member_current"},
		prevStatus: Committed,
		nextStatus: Committed,
		customLog:  idAndDescriptionLog,
	},
	Reject: {
		action:     Reject,
		role:       []string{"admin"},
		prevStatus: Committed,
		nextStatus: Accepted,
		customLog:  idLog,
	},
	Close: {
		action:     Close,
		role:       []string{"admin"},
		prevStatus: Committed,
		nextStatus: Closed,
		customLog: func(eh *eventActionHandler) model.EventLog {
			eh.event.ClosedBy = eh.actor.Id
			return eh.createEventLog(createEventLogArgs{
				Id: eh.actor.Id,
			})
		},
	},
	Update: {
		action:    Update,
		role:      []string{"admin"},
		customLog: idLog,
	},
}

func MakeEventActionHandler(action Action, event *model.Event, identity model.Identity) *eventActionHandler {
	ans := &eventActionHandler{
		action:     eventActionMap[action].action,
		role:       eventActionMap[action].role,
		prevStatus: eventActionMap[action].prevStatus,
		nextStatus: eventActionMap[action].nextStatus,
		customLog:  eventActionMap[action].customLog,
		event:      event,
		actor:      identity,
	}
	return ans
}

// check if the action is valid
func (eh *eventActionHandler) ValidateAction() error {
	roles := []string{eh.actor.Role}
	if eh.actor.Id == eh.event.MemberId {
		roles = append(roles, "member_current")
	}
	clientId, _ := strconv.ParseInt(eh.actor.Id, 10, 64)
	if clientId == eh.event.ClientId {
		roles = append(roles, "client_current")
	}
	if len(eh.role) != 0 {
		exist := false
		for _, role := range roles {
			for _, r := range eh.role {
				if role == r {
					exist = true
					break
				}
			}
		}
		if !exist {
			return MakeServiceError(http.StatusUnprocessableEntity).
				SetMessage("invalid role")
		}
	}
	if eh.prevStatus != eh.event.Status {
		return MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("action not allowed")
	}
	return nil
}

type createEventLogArgs struct {
	Id          string
	Description string
}

func (eh *eventActionHandler) createEventLog(args createEventLogArgs) model.EventLog {
	return model.EventLog{
		EventId:     eh.event.EventId,
		Action:      string(eh.action),
		MemberId:    args.Id,
		Description: args.Description,
		GmtCreate:   GetDate(),
	}
}

func (eh *eventActionHandler) Handle() model.EventLog {
	// set the next status
	eh.event.Status = eh.nextStatus
	var eventLog model.EventLog
	// create log
	if eh.customLog != nil {
		eventLog = eh.customLog(eh)
	} else {
		eventLog = eh.createEventLog(createEventLogArgs{})
	}

	Logger.WithFields(logrus.Fields{
		"event_id":    eventLog.EventId,
		"action":      eventLog.Action,
		"member_id":   eventLog.MemberId,
		"gmt_create":  eventLog.GmtCreate,
		"description": eventLog.Description,
	}).Info("new event action")

	if producer := GetNSQProducer(); producer != nil {
		mapEventLog := map[string]interface{}{
			"event_id":    eventLog.EventId,
			"member_id":   eventLog.MemberId,
			"action":      eventLog.Action,
			"problem":     eh.event.Problem,
			"model":       eh.event.Model,
			"gmt_create":  eventLog.GmtCreate,
			"description": eventLog.Description,
		}
		if eh.actor.Member.Alias != "" {
			mapEventLog["member_alias"] = eh.actor.Member.Alias
		} else {
			mapEventLog["member_alias"] = ""
		}
		jsonMap, _ := json.Marshal(mapEventLog)
		producer.PublishAsync(EventTopic, jsonMap, nil)
	}

	return eventLog
}

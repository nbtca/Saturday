package util

import (
	"net/http"
	"saturday/model"
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

type CustomLogFunc func(*EventActionHandler) model.EventLog

type EventActionHandler struct {
	event       *model.Event
	actor       model.Identity
	action      Action
	role        []string
	prevStatus  string
	nextStatus  string
	Description string
	customLog   CustomLogFunc
}

var idLog CustomLogFunc = func(eh *EventActionHandler) model.EventLog {
	return eh.createEventLog(createEventLogArgs{
		Id: eh.actor.Id,
	})
}

var idAndDescriptionLog CustomLogFunc = func(eh *EventActionHandler) model.EventLog {
	return eh.createEventLog(createEventLogArgs{
		Id:          eh.actor.Id,
		Description: eh.Description,
	})
}

var EventActionMap map[Action]EventActionHandler = map[Action]EventActionHandler{
	Create: {
		action:     Create,
		role:       []string{"client"},
		prevStatus: "",
		nextStatus: Open,
	},
	Accept: {
		action:     Accept,
		role:       []string{"member"},
		prevStatus: Open,
		nextStatus: Accepted,
		customLog: func(eh *EventActionHandler) model.EventLog {
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
		customLog: func(eh *EventActionHandler) model.EventLog {
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
		prevStatus: Accepted,
		nextStatus: Cancelled,
		customLog:  idLog,
	},
	Close: {
		action:     Close,
		role:       []string{"admin"},
		prevStatus: Committed,
		nextStatus: Closed,
		customLog: func(eh *EventActionHandler) model.EventLog {
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

func MakeEventActionHandler(action Action, event *model.Event, identity model.Identity) *EventActionHandler {
	ans := &EventActionHandler{
		action:     EventActionMap[action].action,
		role:       EventActionMap[action].role,
		prevStatus: EventActionMap[action].prevStatus,
		nextStatus: EventActionMap[action].nextStatus,
		customLog:  EventActionMap[action].customLog,
		event:      event,
		actor:      identity,
	}
	return ans
}

// inject the event and actor to the handler
// func (eh *EventActionHandler) Init(action Action, event *model.Event, identity model.Identity) {
// 	eh.action = EventActionMap[action].action
// 	eh.role = EventActionMap[action].role
// 	eh.prevStatus = EventActionMap[action].prevStatus
// 	eh.nextStatus = EventActionMap[action].nextStatus
// 	eh.customLog = EventActionMap[action].customLog
// 	eh.event = event
// 	eh.actor = identity
// }

// check if the action is valid
func (eh *EventActionHandler) ValidateAction() error {
	if len(eh.role) != 0 {
		exist := false
		for _, role := range eh.role {
			if role == eh.actor.Role {
				exist = true
				break
			}
		}
		if !exist {
			return MakeServiceError(http.StatusUnprocessableEntity).
				SetMessage("invalid role")
		}
	}
	if eh.prevStatus != "" && eh.prevStatus != eh.event.Status {
		return MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("action not allowed")
	}
	return nil
}

type createEventLogArgs struct {
	Id          string
	Description string
}

func (eh *EventActionHandler) createEventLog(args createEventLogArgs) model.EventLog {
	return model.EventLog{
		EventId:     eh.event.EventId,
		Action:      string(eh.action),
		MemberId:    args.Id,
		Description: args.Description,
		GmtCreate:   GetDate(),
	}
}

func (eh *EventActionHandler) Handle() model.EventLog {
	// set the next status
	eh.event.Status = eh.nextStatus
	var eventLog model.EventLog
	// create log
	if eh.customLog != nil {
		eventLog = eh.customLog(eh)
	} else {
		eventLog = eh.createEventLog(createEventLogArgs{})
	}
	return eventLog
}

// func (eh *EventActionHandler) Handle() error {
// 	// set the next status
// 	log := eh.operate()
// 	// persist event
// 	err := repo.UpdateEvent(eh.event, &log)
// 	// append log
// 	eh.event.Logs = append(eh.event.Logs, log)
// 	return err
// }

/*
 this function validates the action and then perform action to the event.
 it also persists the event and event log.
*/
// func PerformEventAction(event *model.Event, identity model.Identity, action Action, description ...string) error {
// 	handler := EventActionMap[action]
// 	log.Println(action)
// 	handler.Init(event, identity)
// 	for _, d := range description {
// 		handler.description = fmt.Sprint(handler.description, d)
// 	}
// 	if err := handler.validateAction(); err != nil {
// 		return err
// 	}
// 	return handler.Handle()
// }

package action

import (
	"fmt"
	"log"
	"net/http"
	"saturday/model"
	"saturday/repo"
	"saturday/util"
)

type Identity struct {
	Id   string
	Role string
}

type CustomLogFunc func(*EventActionHandler) model.EventLog

type EventActionHandler struct {
	Event       *model.Event
	Actor       Identity
	action      Action
	role        []string
	prevStatus  string
	nextStatus  string
	description string
	customLog   CustomLogFunc
}

// inject the event and actor to the handler
func (eh *EventActionHandler) Init(event *model.Event, identity Identity) {
	eh.Event = event
	eh.Actor = identity
}

// check if the action is valid
func (eh *EventActionHandler) validateAction() error {
	if len(eh.role) != 0 {
		exist := false
		for _, role := range eh.role {
			if role == eh.Actor.Role {
				exist = true
				break
			}
		}
		if !exist {
			return util.MakeServiceError(http.StatusUnprocessableEntity).
				SetMessage("invalid role")
		}
	}
	if eh.prevStatus != "" && eh.prevStatus != eh.Event.Status {
		log.Println(eh.prevStatus, eh.Event.Status)
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("action not allowed")
	}
	return nil
}

type createEventLogArgs struct {
	Id          string
	Description string
}

func (eh *EventActionHandler) CreateEventLog(args createEventLogArgs) model.EventLog {
	return model.EventLog{
		EventId:     eh.Event.EventId,
		Action:      string(eh.action),
		MemberId:    args.Id,
		Description: args.Description,
		GmtCreate:   util.GetDate(),
	}
}

func (eh *EventActionHandler) Handle() error {
	// set the next status
	eh.Event.Status = eh.nextStatus
	var log model.EventLog
	// create log
	if eh.customLog != nil {
		log = eh.customLog(eh)
	} else {
		log = eh.CreateEventLog(createEventLogArgs{})
	}
	// append log
	// eh.Event.Logs = append(eh.Event.Logs, log)

	// persist event
	err := repo.UpdateEvent(eh.Event, &log)
	eh.Event.Logs = append(eh.Event.Logs, log)
	return err
}

/*
 this function validates the action and then perform action to the event.
 it also persists the event and event log.
*/
func PerformEventAction(event *model.Event, identity Identity, action Action, description ...string) error {
	handler := EventActionMap[action]
	log.Println(action)
	handler.Init(event, identity)
	for _, d := range description {
		handler.description = fmt.Sprint(handler.description, d)
	}
	if err := handler.validateAction(); err != nil {
		return err
	}
	return handler.Handle()
}

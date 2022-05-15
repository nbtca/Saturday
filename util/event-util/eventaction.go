package util

import (
	"fmt"
	"log"
	"net/http"
	"saturday/model"
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
	}
}

func (eh *EventActionHandler) Handle() model.EventLog {
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
	eh.Event.Logs = append(eh.Event.Logs, log)
	return log
}

var idLog CustomLogFunc = func(eh *EventActionHandler) model.EventLog {
	return eh.CreateEventLog(createEventLogArgs{
		Id: eh.Actor.Id,
	})
}

var idAndDescriptionLog CustomLogFunc = func(eh *EventActionHandler) model.EventLog {
	return eh.CreateEventLog(createEventLogArgs{
		Id:          eh.Actor.Id,
		Description: eh.description,
	})
}

/*
 this function validates the action and then handles the event,
 it does not persist the event, only modifies the event.
 you need to call repo methods to persist the event and log.
 description arg corresponses to the EventLog's description field,
 default not to use.
*/
func PerformEventAction(event *model.Event, identity Identity, action Action, description ...string) (model.EventLog, error) {
	handler := EventActionMap[action]
	handler.Init(event, identity)
	handler.Event = event
	handler.Actor = identity
	for _, d := range description {
		handler.description = fmt.Sprint(handler.description, d)
	}
	if err := handler.validateAction(); err != nil {
		return model.EventLog{}, err
	}
	return handler.Handle(), nil
}

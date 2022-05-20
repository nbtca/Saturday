package action

import "saturday/model"

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
)

var idLog CustomLogFunc = func(eh *EventActionHandler) model.EventLog {
	return eh.CreateEventLog(createEventLogArgs{
		Id: eh.Actor.Id,
	})
}

var idAndDescriptionLog CustomLogFunc = func(eh *EventActionHandler) model.EventLog {
	return eh.CreateEventLog(createEventLogArgs{
		Id:          eh.Actor.Id,
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
			eh.Event.MemberId = eh.Actor.Id
			return eh.CreateEventLog(createEventLogArgs{
				Id: eh.Actor.Id,
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
			eh.Event.MemberId = ""
			return eh.CreateEventLog(createEventLogArgs{
				Id: eh.Actor.Id,
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
			eh.Event.ClosedBy = eh.Actor.Id
			return eh.CreateEventLog(createEventLogArgs{
				Id: eh.Actor.Id,
			})
		},
	},
}

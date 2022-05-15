package util

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
		customLog:  idLog,
	},
	Cancel: {
		action:     Cancel,
		role:       []string{"currentClient"},
		prevStatus: Open,
		nextStatus: Cancelled,
		customLog:  idLog,
	},
	Drop: {
		action:     Drop,
		role:       []string{"currentMember"},
		prevStatus: Accepted,
		nextStatus: Cancelled,
		customLog:  idLog,
	},
	Commit: {
		action:     Commit,
		role:       []string{"currentMember"},
		prevStatus: Accepted,
		nextStatus: Committed,
		customLog:  idAndDescriptionLog,
	},
	AlterCommit: {
		action:     AlterCommit,
		role:       []string{"currentMember"},
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
		customLog:  idLog,
	},
}

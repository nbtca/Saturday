package model

type Event struct {
	ClientId          int64  `json:"client_id" db:"client_id"`
	EventId           int64  `json:"event_id" db:"event_id"`
	Model             string `json:"model"`
	Phone             string `json:"phone"`
	Qq                string `json:"qq"`
	ContactPreference string `json:"contact_preference" db:"contact_preference" `
	Problem           string `json:"problem" db:"event_description"`
	// RepairDescription string     `json:"repair_description"`
	MemberId    string     `json:"member_id" db:"member_id"`
	ClosedBy    string     `json:"closed_by" db:"closed_by"`
	Status      string     `json:"status"`
	Logs        []EventLog `json:"logs"`
	GmtCreate   string     `json:"gmt_create" db:"gmt_create"`
	GmtModified string     `json:"gmt_modified" db:"gmt_modified"`
}

type Status struct {
	StatusId int64  `json:"status_id"`
	Status   string `json:"status"`
}
type EventEventStatusRelation struct {
	EventStatusId int64 `json:"event_status_id"`
	EventId       int64 `json:"event_id"`
}

type EventLog struct {
	EventLogId  int64  `json:"log_id"`
	EventId     int64  `json:"-"`
	Description string `json:"description"`
	GmtCreate   string `json:"gmt_create"`
	MemberId    string `json:"member_id"`
	Action      string `json:"action"`
}

type EventActionRelation struct {
	EventLogId    int64 `json:"event_log_id"`
	EventActionId int64 `json:"event_action_id"`
}

type EventAction struct {
	EventActionId int64  `json:"event_action_id"`
	Action        string `json:"action"`
}

package model

type Event struct {
	ClientId          int64  `json:"client_id"`
	EventId           int64  `json:"event_id"`
	Model             string `json:"model"`
	Phone             string `json:"phone"`
	Qq                string `json:"qq"`
	ContactPreference string `json:"contact_preference"`
	EventDescription  string `json:"event_description"`
	RepairDescription string `json:"repair_description"`
	MemberId          string `json:"member_id"`
	ClosedBy          string `json:"closed_by"`
	GmtCreate         string `json:"gmt_create"`
	GmtModified       string `json:"gmt_modified"`
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
	EventLogId  int64  `json:"event_log_id"`
	EventId     int64  `json:"event_id"`
	Description string `json:"description"`
	GmtCreate   string `json:"gmt_create"`
	MemberId    string `json:"member_id"`
}
type EventActionRelation struct {
	EventLogId    int64 `json:"event_log_id"`
	EventActionId int64 `json:"event_action_id"`
}

package model

type EventLog struct {
	EventLogId int64 `json:"event_log_id"`
	EventId int64 `json:"event_id"`
	Description string `json:"description"`
	GmtCreate string `json:"gmt_create"`
	MemberId string `json:"member_id"`
}
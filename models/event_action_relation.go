package model

type EventActionRelation struct {
	EventLogId int64 `json:"event_log_id"`
	EventActionId int64 `json:"event_action_id"`
}
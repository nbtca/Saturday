package dto

type EventID struct {
	EventID int64 `uri:"EventId" json:"event_id" binding:"required"`
}

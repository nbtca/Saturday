package dto

type EventID struct {
	EventID int64 `uri:"EventId" json:"event_id" binding:"required"`
}
type CommitReq struct {
	Content string `json:"content" binding:"required"`
}

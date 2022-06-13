package dto

type EventID struct {
	EventID int64 `uri:"EventId" json:"event_id" binding:"required"`
}

type CommitReq struct {
	Content string `json:"content"`
}

type AlterCommitReq struct {
	Content string `json:"content"`
}

type UpdateReq struct {
	Phone   string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ      string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
	Problem string `json:"problem" db:"problem" binding:"omitempty,max=1000"`
}
type CreateEventReq struct {
	ClientId          int64  `json:"client_id" db:"client_id"`
	Model             string `json:"model" binding:"omitempty,max=40"`
	Phone             string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ                string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
	ContactPreference string `json:"contact_preference" db:"contact_preference" `
	Problem           string `json:"problem" db:"problem" binding:"omitempty,max=1000"`
}

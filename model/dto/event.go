package dto

type EventID struct {
	EventID int64 `uri:"EventId" json:"eventId" binding:"required"`
}

type CommitRequest struct {
	Content string `json:"content"`
	Size  string `json:"size"`
}

type AlterCommitRequest struct {
	Content string `json:"content"`
	Size  string `json:"size"`
}

type UpdateRequest struct {
	Phone             string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ                string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
	Problem           string `json:"problem" db:"problem" binding:"omitempty,max=1000"`
	Model             string `json:"model" binding:"omitempty,max=40"`
	ContactPreference string `json:"contactPreference" db:"contact_preference" `
	Size              string `json:"size" db:"size" `
}
type CreateEventRequest struct {
	ClientId          int64  `json:"clientId" db:"client_id"`
	Model             string `json:"model" binding:"omitempty,max=40"`
	Phone             string `json:"phone" binding:"required,omitempty,len=11,numeric"`
	QQ                string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
	ContactPreference string `json:"contactPreference" db:"contact_preference" `
	Problem           string `json:"problem" db:"problem" binding:"required,omitempty,max=1000"`
}

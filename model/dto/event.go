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
	Phone   string `json:"phone" binding:"omitempty,len=11"`
	QQ      string `json:"qq" binding:"omitempty,min=5,max=9"`
	Problem string `json:"problem" db:"problem"`
}

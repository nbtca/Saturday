package model

type MemberStatus struct {
	MemberStatusId int64 `json:"member_status_id" db:"member_status_id"`
	Role string `json:"role" db:"role"`
}
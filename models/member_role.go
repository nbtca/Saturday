package model

type MemberRole struct {
	MemberRoleId int64 `json:"member_role_id" db:"member_role_id"`
	Role string `json:"role" db:"role"`
}
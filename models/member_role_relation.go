package model

type MemberRoleRelation struct {
	MemberId string `json:"member_id"`
	RoleId int64 `json:"role_id"`
}
package dto

import model "saturday/src/model"

type MemberId struct {
	MemberId string `uri:"MemberId" json:"member_id" binding:"required,len=10,numeric"`
}

type CreateMemberReq struct {
	MemberId
	Alias   string `json:"alias"`
	Name    string `json:"name" binding:"required,min=2,max=4"`
	Section string `json:"section" binding:"required,section"`
	Profile string `json:"profile"`
	Phone   string `json:"phone" binding:"omitempty,len=11,numeric"`
	Qq      string `json:"qq" binding:"omitempty,min=5,max=12,numeric"`
}

type CreateMemberTokenReq struct {
	MemberId string `uri:"MemberId" binding:"required,len=10,numeric"`
	Password string `json:"password" binding:"required"`
}

type CreateMemberTokenResponse struct {
	model.Member
	Token string `json:"token"`
}

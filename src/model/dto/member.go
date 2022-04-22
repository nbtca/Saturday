package dto

import model "gin-example/src/model"

type CreateMemberTokenReq struct {
	MemberId string `json:"member_id" validate:"required,len=10,numeric"`
	Password string `json:"password" validate:"required"`
}

type CreateMemberTokenResponse struct {
	model.Member
	Token string `json:"token"`
}

type Page struct {
	Offset uint64 `json:"-" validate:"min=0"`
	Limit  uint64 `json:"-" validate:"min=0"`
}

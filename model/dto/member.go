package dto

import model "github.com/nbtca/saturday/model"

type MemberId struct {
	MemberId string `uri:"MemberId" json:"memberId" binding:"required,len=10,numeric"`
}

type CreateMemberRequest struct {
	MemberId string `uri:"MemberId" json:"memberId" binding:"required,len=10,numeric"`
	LogtoId  string `json:"logtoId" binding:"omitempty"`
	Name     string `json:"name" binding:"required,min=2,max=4"`
	Section  string `json:"section" binding:"required,section"`
	Alias    string `json:"alias" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	Profile  string `json:"profile" binding:"omitempty,max=1000"`
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ       string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
	Role     string `json:"role" binding:"required"`
}

type CreateMemberWithLogtoRequest struct {
	MemberId string `uri:"MemberId" json:"memberId" binding:"required,len=10,numeric"`
	LogtoId  string `json:"logtoId" binding:"omitempty"`
	Name     string `json:"name" binding:"required,min=2,max=4"`
	Section  string `json:"section" binding:"required,section"`
	Alias    string `json:"alias" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	Profile  string `json:"profile" binding:"omitempty,max=1000"`
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ       string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
}

type CreateMemberTokenRequest struct {
	MemberId string `uri:"MemberId" binding:"required,len=10,numeric"`
	Password string `json:"password" binding:""`
}

type CreateMemberTokenResponse struct {
	model.Member
	Token string `json:"token"`
}

type UpdateMemberBasicRequest struct {
	MemberId string `uri:"MemberId" json:"memberId" binding:"required,len=10,numeric"`
	Name     string `json:"name" binding:"omitempty,min=2,max=4"`
	Section  string `json:"section" binding:"omitempty,section"`
	Role     string `json:"role" binding:"omitempty"`
}

type UpdateMemberRequest struct {
	MemberId string
	Alias    string `json:"alias" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	Profile  string `json:"profile" binding:"omitempty,max=1000"`
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ       string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
	Password string `json:"password" binding:"omitempty,max=20"`
}

type ActivateMemberRequest struct {
	MemberId string
	Password string `json:"password" binding:"required,max=50"`
	Alias    string `json:"alias" binding:"omitempty,max=20"`
	Profile  string `json:"profile" binding:"omitempty,max=1000"`
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
	QQ       string `json:"qq" binding:"omitempty,min=5,max=20,numeric"`
}

type UpdateAvatarRequest struct {
	Url string `json:"url" binding:"required"`
}

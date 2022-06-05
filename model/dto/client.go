package dto

import "saturday/model"

type WxLoginReq struct {
	Code string `json:"code" binding:"required"`
}

type ClientTokenResponse struct {
	Token string `json:"token"`
	model.Client
}

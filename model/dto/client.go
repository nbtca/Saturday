package dto

import "github.com/nbtca/saturday/model"

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type ClientTokenResponse struct {
	Token string `json:"token"`
	model.Client
}

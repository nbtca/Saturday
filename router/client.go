package router

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

type ClientRouter struct{}

func (ClientRouter) CreateTokenViaWeChat(c context.Context, input *struct {
	Body struct {
		Code string `json:"code"`
	}
}) (*util.CommonResponse[dto.ClientTokenResponse], error) {
	openid, err := util.CodeToSession(input.Body.Code)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	client, err := service.ClientServiceApp.GetClientByOpenId(openid)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if client == (model.Client{}) {
		client, err = service.ClientServiceApp.CreateClientByOpenId(openid)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
	}
	token, err := service.ClientServiceApp.CreateTokenViaWechat(client)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(dto.ClientTokenResponse{
		Token:  token,
		Client: client,
	}), nil
}

var ClientRouterApp = ClientRouter{}

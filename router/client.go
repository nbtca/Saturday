package router

import (
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type ClientRouter struct{}

func (ClientRouter) CreateTokenViaWeChat(c *gin.Context) {
	wxLoginRequest := &dto.WxLoginRequest{}
	if err := util.BindAll(c, wxLoginRequest); util.CheckError(c, err) {
		return
	}
	openid, err := util.CodeToSession(wxLoginRequest.Code)
	if util.CheckError(c, err) {
		return
	}
	client, err := service.ClientServiceApp.GetClientByOpenId(openid)
	if util.CheckError(c, err) {
		return
	}
	if client == (model.Client{}) {
		client, err = service.ClientServiceApp.CreateClientByOpenId(openid)
		if util.CheckError(c, err) {
			return
		}
	}
	token, err := service.ClientServiceApp.CreateTokenViaWechat(client)
	if util.CheckError(c, err) {
		return
	}
	res := dto.ClientTokenResponse{
		Token:  token,
		Client: client,
	}
	c.JSON(200, res)
}

var ClientRouterApp = ClientRouter{}

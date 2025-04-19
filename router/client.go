package router

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/nbtca/saturday/middleware"
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
	token, err := service.ClientServiceApp.CreateClientToken(client)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(dto.ClientTokenResponse{
		Token:  token,
		Client: client,
	}), nil
}

func (ClientRouter) CreateTokenViaLogto(c *gin.Context) {
	user := c.Value("user").(middleware.AuthContextUser)
	logtoId := user.UserInfo.Sub
	if logtoId == "" {
		c.Error(huma.Error422UnprocessableEntity("user not found"))
		// return nil, huma.Error422UnprocessableEntity("user not found")
	}

	client, err := service.ClientServiceApp.CreateClientByLogtoIdIfNotExists(logtoId)
	if err != nil {
		c.Error(huma.Error422UnprocessableEntity(err.Error()))
	}

	token, err := service.ClientServiceApp.CreateClientToken(client)
	if err != nil {
		// return nil, huma.Error422UnprocessableEntity(err.Error())
		c.Error(huma.Error422UnprocessableEntity(err.Error()))
	}
	response := dto.ClientTokenResponse{
		Token:  token,
		Client: client,
	}
	c.JSON(200, response)
}

var ClientRouterApp = ClientRouter{}

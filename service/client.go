package service

import (
	"fmt"
	"saturday/model"
	"saturday/repo"
	"saturday/util"
)

type ClientService struct {
}

func (service ClientService) CreateTokenViaWechat(client model.Client) (string, error) {
	res, err := util.CreateToken(util.Payload{Who: fmt.Sprint(client.ClientId), Role: "client"})
	return res, err
}

func (service ClientService) GetClientByOpenId(openId string) (model.Client, error) {
	client, err := repo.GetClientByOpenId(openId)
	if err != nil {
		return model.Client{}, err
	}
	return client, nil
}

func (service ClientService) CreateClientByOpenId(openId string) (model.Client, error) {
	client := model.Client{OpenId: openId}
	err := repo.CreateClient(&client)
	if err != nil {
		return model.Client{}, err
	}
	return client, nil
}

var ClientServiceApp = ClientService{}
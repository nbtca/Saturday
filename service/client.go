package service

import (
	"fmt"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"
)

type ClientService struct {
}

func (service ClientService) CreateClientToken(client model.Client) (string, error) {
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

func (service ClientService) GetClientByLogtoId(logtoId string) (model.Client, error) {
	client, err := repo.GetClientByLogtoId(logtoId)
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

func (service ClientService) CreateClientByLogtoId(logtoId string) (model.Client, error) {
	client := model.Client{LogtoId: logtoId}
	err := repo.CreateClient(&client)
	if err != nil {
		return model.Client{}, err
	}
	return client, nil
}

func (service ClientService) CreateClientByLogtoIdIfNotExists(logtoId string) (model.Client, error) {
	client, err := service.GetClientByLogtoId(logtoId)
	if err != nil {
		return model.Client{}, err
	}
	if client == (model.Client{}) {
		client, err = service.CreateClientByLogtoId(logtoId)
		if err != nil {
			return model.Client{}, err
		}
	}
	return client, nil
}

var ClientServiceApp = ClientService{}

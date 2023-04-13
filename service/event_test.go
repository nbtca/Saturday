package service_test

import (
	"log"
	"saturday/model"
	"saturday/service"
	"testing"
)

func TestEventService_SendActionNotificationViaPushDeer(t *testing.T) {
	service := service.EventService{}
	err := service.SendActionNotifyViaPushDeer(&model.Event{
		Model:     "model",
		Problem:   "problem",
		GmtCreate: "gmtCreate",
	}, "test")
	if err != nil {
		log.Print(err)
	}
}

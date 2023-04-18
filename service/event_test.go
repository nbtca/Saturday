package service_test

import (
	"log"
	"testing"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/service"
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

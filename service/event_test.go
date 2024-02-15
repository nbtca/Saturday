package service_test

import (
	"log"
	"os"
	"testing"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/service"
)

func TestSendActionNotifyViaRPC(t *testing.T) {
	os.Setenv("RPC_ADDRESS", ":8000")
	service := service.EventService{}
	err := service.SendActionNotifyViaRPC(&model.EventActionNotifyRequest{
		Model:     "model",
		Problem:   "problem",
		GmtCreate: "gmtCreate",
	})
	if err != nil {
		log.Print(err)
	}
}

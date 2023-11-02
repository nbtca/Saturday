package service_test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

func TestMain(m *testing.M) {

	if err := godotenv.Load("../.env"); err != nil {
		util.Logger.Warning("Error loading .env file")
		log.Println(err)
	}
	m.Run()
}

func TestFetchLogtoToken(t *testing.T) {
	res, err := service.LogtoServiceApp.FetchLogtoToken("https://default.logto.app/api", "all")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestFetchLogtoUser(t *testing.T) {
	res, err := service.LogtoServiceApp.FetchLogtoToken()
	if err != nil {
		t.Error(err)
	}
	token := res["access_token"].(string)
	userId := "chmz1itz83qq"
	user, err := service.LogtoServiceApp.FetchUserById(userId, "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

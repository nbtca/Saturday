package service_test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {

	if err := godotenv.Load("../.env"); err != nil {
		util.Logger.Warning("Error loading .env file")
		log.Println(err)
	}
	m.Run()
}

func TestFetchLogtoToken(t *testing.T) {
	service.LogtoServiceApp = service.MakeLogtoService(viper.GetString("logto.endpoint"))
	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestFetchLogtoUser(t *testing.T) {
	service.LogtoServiceApp = service.MakeLogtoService(viper.GetString("logto.endpoint"))
	userId := viper.GetString("TESTING_LOGTO_USER_ID")
	user, err := service.LogtoServiceApp.FetchUserById(userId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestFetchLogtoUsers(t *testing.T) {
	service.LogtoServiceApp = service.MakeLogtoService(viper.GetString("logto.endpoint"))
	// userId := viper.GetString("TESTING_LOGTO_USER_ID")
	user, err := service.LogtoServiceApp.FetchUsers(service.FetchLogtoUsersRequest{
		Page:         1,
		PageSize:     10,
		SearchParams: map[string]interface{}{"search.primaryEmail": "clas.wen@icloud.com"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user[0].Identities["github"].UserId)
}

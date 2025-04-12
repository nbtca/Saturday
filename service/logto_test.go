package service_test

import (
	"log"
	"os"
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
	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))
	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestFetchLogtoUser(t *testing.T) {
	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))
	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if err != nil {
		t.Error(err)
	}
	token := res["access_token"].(string)
	userId := os.Getenv("LOGTO_TEST_USER_ID")
	user, err := service.LogtoServiceApp.FetchUserById(userId, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestFetchLogtoUsers(t *testing.T) {
	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))
	// userId := os.Getenv("LOGTO_TEST_USER_ID")
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

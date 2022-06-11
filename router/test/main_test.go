package router_test

import (
	"log"
	"saturday/repo"
	"saturday/router"
	"saturday/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var r *gin.Engine
var db *sqlx.DB
var mockDB *util.MockDB

func TestMain(m *testing.M) {
	util.InitValidator()

	mockDB = util.MakeMockDB("../../assets")
	db, err := mockDB.Start()
	if err != nil {
		log.Fatal(err)
	}
	repo.SetDB(db)
	defer mockDB.Close()

	r = router.SetupRouter()
	m.Run()
}

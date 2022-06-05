package router_test

import (
	"saturday/repo"
	"saturday/router"
	"saturday/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var r *gin.Engine
var db *sqlx.DB

func TestMain(m *testing.M) {
	util.InitValidator()
	db, _ = util.GetDB()
	repo.SetDB(db)
	defer repo.CloseDB()
	defer util.CloseResource()

	r = router.SetupRouter()

	m.Run()
}

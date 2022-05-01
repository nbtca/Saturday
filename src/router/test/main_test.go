package router_test

import (
	"saturday/src/repo"
	"saturday/src/router"
	"saturday/src/util"
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
	defer util.Close()

	r = router.SetupRouter()

	m.Run()
}

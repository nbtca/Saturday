package apitest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"saturday/src/repo"
	"saturday/src/router"
	"saturday/src/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	os.Chdir("./../..")
	util.InitEnv()

	repo.InitDB()
	defer repo.CloseDB()

	r = router.SetupRouter()

	m.Run()
}

func TestPingRoute(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

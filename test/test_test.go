package test

import (
	"encoding/json"
	"gin-example/src/repo"
	"gin-example/src/router"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestGetResourceById(t *testing.T) {
	os.Chdir("./..")
	r := router.SetupRouter()
	repo.InitDB()
	defer repo.CloseDB()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/2333333333", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code) // or what value you need it to be

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, gin.H{
		"member_id":    "2333333333",
		"alias":        "huaji",
		"password":     "123456",
		"name":         "滑稽",
		"section":      "计算机233",
		"role":         "member",
		"profile":      "",
		"phone":        "",
		"qq":           "",
		"avatar":       "",
		"created_by":   "",
		"gmt_create":   "2022-04-23 15:49:59",
		"gmt_modified": "2022-04-23 15:50:01",
	}, got) // want is a gin.H that contains the wanted map.
}

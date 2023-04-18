package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

func TestAuth(t *testing.T) {
	rawCase, err := util.GetCsvMap("testdata/auth.csv")
	if err != nil {
		t.Error(err)
	}
	for _, rc := range rawCase {
		t.Run(rc["CaseId"], func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}
			auth := ""
			if rc["role"] == "T" {
				auth = "member"
			}
			if rc["valid"] != "T" {
				auth = "INVALID"
			}
			if rc["authorization"] == "T" {
				token := util.GenToken(auth, "2333333333")
				c.Request.Header.Add("Authorization", token)
			}
			middleware.Auth("member")(c)

			if rc["error"] == "T" && w.Code != 422 {
				t.Error("error case should return 422")
			}

			if rc["success"] == "T" {
				if w.Code != 200 {
					t.Error("success case should return 200")
				}
				r, exist := c.Get("role")
				if !exist {
					t.Error("role should be set")
				}
				if r != "member" {
					t.Error("role should be member")
				}
				id, exist := c.Get("id")
				if !exist {
					t.Error("id should be set")
				}
				if id != "2333333333" {
					t.Error("id should be 2333333333")
				}
			}
		})
	}
}

package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TODO assert
func DataHandler(data TestData, w *httptest.ResponseRecorder) {
	body := data.RequestBody
	bodyData, _ := json.Marshal(body)
	reader := bytes.NewReader(bodyData)
	req, _ := http.NewRequest(data.Method, data.Url, reader)
	r.ServeHTTP(w, req)
	// assert.Equal(t, 200, w.Code)
}

func TestCreateMember(t *testing.T) {
	w := httptest.NewRecorder()
	for _, data := range CreateMemberData {
		var reader *bytes.Reader
		body := data.RequestBody
		if body != nil {
			bodyData, _ := json.Marshal(body)
			reader = bytes.NewReader(bodyData)
		}
		req, _ := http.NewRequest(data.Method, data.Url, reader)
		r.ServeHTTP(w, req)
		var got gin.H
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, data.Code, w.Code)
		// TODO check body
	}
}

func TestGetMemberById(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/members/2333333333", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code) // or what value you need it to be

	// var got gin.H
	// err := json.Unmarshal(w.Body.Bytes(), &got)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// assert.Equal(t, gin.H{
	// 	"member_id":    "2333333333",
	// 	"alias":        "huaji",
	// 	"password":     "123456",
	// 	"name":         "滑稽",
	// 	"section":      "计算机233",
	// 	"role":         "member",
	// 	"profile":      "",
	// 	"phone":        "",
	// 	"qq":           "",
	// 	"avatar":       "",
	// 	"created_by":   "",
	// 	"gmt_create":   "2022-04-23 15:49:59",
	// 	"gmt_modified": "2022-04-23 15:50:01",
	// }, got) // want is a gin.H that contains the wanted map.
}

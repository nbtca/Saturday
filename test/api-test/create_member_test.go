package apitest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateMember(t *testing.T) {
	w := httptest.NewRecorder()
	body := gin.H{
		"alias":    "huaji",
		"password": "123456",
		"section":  "计算机233",
	}
	data, _ := json.Marshal(body)
	reader := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", "/members/2333333000", reader)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code) // or what value you need it to be

}

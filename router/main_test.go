package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"saturday/repo"
	"saturday/router"
	"saturday/util"
	"testing"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var mockDB *util.MockDB

func TestMain(m *testing.M) {
	util.InitValidator()

	mockDB = util.MakeMockDB("../assets")
	db, err := mockDB.Start()
	if err != nil {
		log.Fatal(err)
	}
	repo.SetDB(db)
	defer mockDB.Close()

	r = router.SetupRouter()
	m.Run()
}

type APITestCase struct {
	Name     string
	Request  Request
	Response Response
}

type Request struct {
	Method string
	Url    string
	Auth   string
	Body   gin.H
}
type Response struct {
	Code int
	Body gin.H
}

func (t APITestCase) compare(got gin.H) error {
	if t.Response.Body == nil {
		return nil
	}
	if t.Response.Code != 200 {
		if t.Response.Body["message"] != got["message"] {
			return fmt.Errorf("inconsistent message\n expected:%v\n got:%v", t.Response.Body["message"], got["message"])
		}
		return nil
	}
	if len(got) > len(t.Response.Body) {
		return fmt.Errorf("extra field in response body\n expected:%v \n got:%v", t.Response.Body, got)
	}
	for key := range t.Response.Body {
		v := got[key]
		if v == nil {
			return fmt.Errorf("missing field in response body:%v", key)
		}
		// ignore fields
		if key == "gmt_create" || key == "gmt_modified" {
			continue
		}
		if t.Response.Body[key] == "IGNORE" {
			continue
		}
		if v != t.Response.Body[key] {
			return fmt.Errorf("inconsistent response body\n key:%v\n expected:%v \n got:%v", key, t.Response.Body[key], v)
		}
	}
	return nil
}

func (tc APITestCase) Test() error {
	if err := mockDB.SetSchema(); err != nil {
		return err
	}
	w := httptest.NewRecorder()
	var reader *bytes.Reader
	body := tc.Request.Body
	var req *http.Request
	if body != nil {
		bodyData, _ := json.Marshal(body)
		reader = bytes.NewReader(bodyData)
		req, _ = http.NewRequest(tc.Request.Method, tc.Request.Url, reader)
	} else {
		req, _ = http.NewRequest(tc.Request.Method, tc.Request.Url, nil)
	}
	if tc.Request.Auth != "" {
		req.Header.Add("Authorization", tc.Request.Auth)
	}
	r.ServeHTTP(w, req)
	if tc.Response.Code != w.Code {
		return fmt.Errorf("inconsistent code\n expected:%v\n got:%v", tc.Response.Code, w.Code)
	}
	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		return err
	}
	return tc.compare(got)
}

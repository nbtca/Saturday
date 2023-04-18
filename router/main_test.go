package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var mockDB *util.MockDB

const IGNORE = "IGNORE"

func TestMain(m *testing.M) {
	util.InitValidator()

	mockDB = util.MakeMockDB("../assets")
	db, err := mockDB.Start()
	if err != nil {
		log.Fatal(err)
	}
	repo.SetDB(db)
	// defer mockDB.Close()
	defer mockDB.CloseDb()

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
		if key == "gmtCreate" || key == "gmtModified" {
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

// Run() perform the request without comparing the response
func (tc APITestCase) Run() error {
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
	return nil
}

// Test() perform the request and compare the response
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

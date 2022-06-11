package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"saturday/util"

	"github.com/gin-gonic/gin"
)

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

func GenToken(auth string, id ...string) string {
	if auth == "INVALID" {
		return "Invalid"
	} else if auth == "EXPIRED" {
		return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTM0NDMxNjIsImRhdGEiOnsidWlkIjoiNmYxZjk3MDItNjZkNi00NDdiLThlNTUtNWYwYzY0N2M4ZDNhIiwicm9sZSI6InVzZXIifSwiaWF0IjoxNjUzMzU2NzYyfQ.ocAxJGhw6Xt2vt7bwGcMeRPLOQOmaspznyu9aI7G670"
	} else if auth == "NONE" {
		return ""
	}
	token, _ := util.CreateToken(util.Payload{Who: id[0], Role: auth})
	return token
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
			return fmt.Errorf("inconsistent response body\n expected:%v \n got:%v", t.Response.Body[key], v)
		}
	}
	return nil
}

func (tc APITestCase) Test() error {
	if err := mockDB.SetSchema(db); err != nil {
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
		// token, _ := util.CreateToken(util.Payload{Who: "2333333333", Role: tc.Request.Role})
		req.Header.Add("Authorization", tc.Request.Auth)
	}
	r.ServeHTTP(w, req)
	if tc.Response.Code != w.Code {
		log.Println(body)
		return fmt.Errorf("inconsistent code\n expected:%v\n got:%v", tc.Response.Code, w.Code)
	}
	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		log.Println("json Unmarshal err")
		return err
	}
	return tc.compare(got)
}

// type DetailError struct {
// 	resource string
// 	field    string
// 	error    string
// }

// type FailBody struct {
// 	message string
// 	errors  []DetailError
// }

// func DataHandler(data APITestCase) error {
// 	if err := mockDB.SetSchema(db); err != nil {
// 		return err
// 	}
// 	w := httptest.NewRecorder()
// 	var reader *bytes.Reader
// 	body := tc.Request.Body
// 	var req *http.Request
// 	if body != nil {
// 		bodyData, _ := json.Marshal(body)
// 		reader = bytes.NewReader(bodyData)
// 		req, _ = http.NewRequest(tc.Request.Method, tc.Request.Url, reader)
// 	} else {
// 		req, _ = http.NewRequest(tc.Request.Method, tc.Request.Url, nil)
// 	}
// 	if tc.Request.Role != "" {
// 		token, _ := util.CreateToken(util.Payload{Who: "2333333333", Role: tc.Request.Role})
// 		req.Header.Add("Authorization", token)
// 	}
// 	r.ServeHTTP(w, req)
// 	if tc.Response.Code != w.Code {
// 		log.Println(body)
// 		return fmt.Errorf("inconsistent code\n expected:%v\n got:%v", tc.Response.Code, w.Code)
// 	}
// 	var got gin.H
// 	err := json.Unmarshal(w.Body.Bytes(), &got)
// 	if err != nil {
// 		log.Println("json Unmarshal err")
// 		return err
// 	}
// 	return tc.compare(got)
// }

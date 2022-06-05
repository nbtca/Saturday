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

// type H map[string]interface{}

type TestData struct {
	Name string
	// Method       string
	// Url          string
	// RequestBody  H
	// Code         int
	// ResponseBody H
	Req Request
	Res Response
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

// type DetailError struct {
// 	resource string
// 	field    string
// 	error    string
// }

// type FailBody struct {
// 	message string
// 	errors  []DetailError
// }

func DataHandler(data TestData) error {
	if err := util.SetSchema(db); err != nil {
		return err
	}
	w := httptest.NewRecorder()
	var reader *bytes.Reader
	body := data.Req.Body
	var req *http.Request
	if body != nil {
		bodyData, _ := json.Marshal(body)
		reader = bytes.NewReader(bodyData)
		req, _ = http.NewRequest(data.Req.Method, data.Req.Url, reader)
	} else {
		req, _ = http.NewRequest(data.Req.Method, data.Req.Url, nil)
	}
	token, _ := util.CreateToken(util.Payload{Who: "2333333333", Role: "member"})
	req.Header.Add("Authorization", token)
	r.ServeHTTP(w, req)
	if data.Res.Code != w.Code {
		log.Println(body)
		return fmt.Errorf("inconsistent code\n expected:%v\n got:%v", data.Res.Code, w.Code)
	}

	// TODO remove this and test error response body
	if data.Res.Code != 200 {
		return nil
	}
	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		log.Println("json Unmarshal err")
		return err
	}
	// no body to check
	if data.Res.Body == nil {
		return nil
	}
	if len(got) > len(data.Res.Body) {
		return fmt.Errorf("extra field in response body\n expected:%v \n got:%v", data.Res.Body, got)
	}
	for key := range data.Res.Body {
		v := got[key]
		if v == nil {
			return fmt.Errorf("missing field in response body:%v", key)
		}
		// ignore fields
		if key == "gmt_create" || key == "gmt_modified" {
			continue
		}
		if v != data.Res.Body[key] {
			return fmt.Errorf("inconsistent response body\n expected:%v \n got:%v", data.Res.Body, got)
		}
	}
	return nil
}

package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Body struct {
	OpenId      string `json:"openid"`
	Session_key string `json:"session_key"`
	Union       string `json:"unionid"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func CodeToSession(code string) (string, error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + os.Getenv("APPID") + "&secret=" + os.Getenv("SECRET") + "&js_code=" + code + "&grant_type=authorization_code"
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	data := &Body{}
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if data.ErrCode != 0 {
		return "", MakeServiceError(http.StatusUnprocessableEntity).SetMessage(data.ErrMsg)
	}
	return data.OpenId, nil
}

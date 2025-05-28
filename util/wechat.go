package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type jscode2sessionBody struct {
	OpenId      string `json:"openid"`
	Session_key string `json:"session_key"`
	Union       string `json:"unionid"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func CodeToSession(code string) (string, error) {
	appid := viper.GetString("wechat.appid")
	secret := viper.GetString("wechat.secret")
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code", appid, secret, code)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	data := &jscode2sessionBody{}
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if data.ErrCode != 0 {
		return "", MakeServiceError(http.StatusUnprocessableEntity).SetMessage(data.ErrMsg)
	}
	return data.OpenId, nil
}

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nbtca/saturday/util"
)

type LogtoService struct {
	BaseURL string
}

func (l LogtoService) FetchLogtoToken(resource string, scope string) (map[string]interface{}, error) {
	tokenURL, _ := url.JoinPath(l.BaseURL, "/oidc/token")

	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("resource", resource)
	params.Add("scope", scope)
	payload := strings.NewReader(params.Encode())

	req, _ := http.NewRequest("POST", tokenURL, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	authString := util.CreateBasicAuth(os.Getenv("LOGTO_APPID"), os.Getenv("LOGTO_APP_SECRET"))
	req.Header.Add("Authorization", "Basic "+authString)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf(string(rawBody))
	}

	return body, nil
}

func (l LogtoService) FetchUserById(userId string, token string) (map[string]interface{}, error) {
	requestURL, _ := url.JoinPath(l.BaseURL, "/api/users/", userId)
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("Authorization", token)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf(string(rawBody))
	}
	return body, nil
}

func MakeLogtoService(endpoint string) LogtoService {
	return LogtoService{
		BaseURL: endpoint,
	}
}

var LogtoServiceApp LogtoService

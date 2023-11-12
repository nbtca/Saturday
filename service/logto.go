package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/util"
)

var DefaultLogtoResource = "https://default.logto.app/api"

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

func (l LogtoService) PatchUserById(userId string, data dto.PatchLogtoUserRequest, token string) (map[string]interface{}, error) {
	requestURL, _ := url.JoinPath(l.BaseURL, "/api/users/", userId)

	var payload bytes.Buffer
	if err := json.NewEncoder(&payload).Encode(data); err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("PATCH", requestURL, &payload)
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

func (l LogtoService) FetchUserByToken(token string, accessToken string) (map[string]interface{}, error) {
	jwksURL, err := url.JoinPath(os.Getenv("LOGTO_ENDPOINT"), "/oidc/jwks")
	if err != nil {
		return nil, err
	}

	invalidTokenError := util.
		MakeServiceError(http.StatusUnprocessableEntity).
		AddDetailError("member", "logto token", "invalid token")
	_, claims, error := util.ParseTokenWithJWKS(jwksURL, token)
	if error != nil {
		return nil, invalidTokenError.SetMessage("Invalid token" + error.Error())
	}
	// check issuer
	expectedIssuer, _ := url.JoinPath(os.Getenv("LOGTO_ENDPOINT"), "/oidc")
	if claims.Issuer != expectedIssuer {
		return nil, invalidTokenError.SetMessage("Invalid token, invalid issuer")
	}
	// check audience
	// TODO move current resource indicator to config
	// expectedAudience := "https://api.nbtca.space/v2"
	// if claims.Audience != expectedAudience {
	// 	c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token").Build())
	// 	return
	// }
	// TODO check scope

	userId := claims.Subject
	user, err := l.FetchUserById(userId, "Bearer "+accessToken)
	if err != nil {
		return nil, invalidTokenError.SetMessage("Invalid token")
	}
	return user, nil
}

func MakeLogtoService(endpoint string) LogtoService {
	return LogtoService{
		BaseURL: endpoint,
	}
}

var LogtoServiceApp LogtoService

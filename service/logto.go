package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/util"
	"github.com/spf13/viper"
)

var DefaultLogtoResource = "https://default.logto.app/api"

type LogtoService struct {
	BaseURL string
}

var logtoToken struct {
	sync.Mutex
	value     string
	expiresAt time.Time
}

func (l LogtoService) getToken() (string, error) {
	logtoToken.Lock()
	defer logtoToken.Unlock()
	if logtoToken.value != "" && time.Now().Before(logtoToken.expiresAt) {
		return logtoToken.value, nil
	}
	res, err := l.FetchLogtoToken(DefaultLogtoResource, "all")
	if err != nil {
		return "", err
	}
	token, _ := res["access_token"].(string)
	expiresIn, _ := res["expires_in"].(float64)
	logtoToken.value = token
	logtoToken.expiresAt = time.Now().Add(time.Duration(expiresIn)*time.Second - time.Minute)
	return token, nil
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
	authString := util.CreateBasicAuth(viper.GetString("logto.appid"), viper.GetString("logto.app_secret"))
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

type FetchLogtoUsersRequest struct {
	Page         int32
	PageSize     int32
	SearchParams map[string]interface{}
}
type LogtoUserIdentities struct {
	UserId  string                 `json:"userId"`
	Details map[string]interface{} `json:"details"`
}

type FetchLogtoUsersResponse struct {
	Id            string                         `json:"id"`
	UserName      string                         `json:"username"`
	PrimaryEmail  string                         `json:"primaryEmail"`
	PrimaryPhone  string                         `json:"primaryPhone"`
	Name          string                         `json:"name"`
	Avatar        string                         `json:"avatar"`
	Identities    map[string]LogtoUserIdentities `json:"identities"`
	CustomerData  map[string]interface{}         `json:"customData"`
	SSOIdentities []map[string]string            `json:"ssoIdentities"`
}

func (l LogtoService) FetchUsers(request FetchLogtoUsersRequest) ([]FetchLogtoUsersResponse, error) {
	if request.Page < 1 {
		request.Page = 1
	}
	if request.PageSize < 1 {
		request.PageSize = 10
	}
	query := url.Values{
		"page":      {fmt.Sprint(request.Page)},
		"page_size": {fmt.Sprint(request.PageSize)},
	}
	for key, value := range request.SearchParams {
		query.Add(key, fmt.Sprint(value))
	}
	requestURL, _ := url.JoinPath(l.BaseURL, "/api/users")
	requestURL = requestURL + "?" + query.Encode()
	log.Println(requestURL)
	req, _ := http.NewRequest("GET", requestURL, nil)
	token, err := l.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body []FetchLogtoUsersResponse
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf(string(rawBody))
	}
	return body, nil

}

func (l LogtoService) FetchUserById(userId string) (*FetchLogtoUsersResponse, error) {
	requestURL, _ := url.JoinPath(l.BaseURL, "/api/users/", userId)
	req, _ := http.NewRequest("GET", requestURL, nil)
	token, err := l.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body FetchLogtoUsersResponse
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf(string(rawBody))
	}
	return &body, nil
}

func (l LogtoService) PatchUserById(userId string, data dto.PatchLogtoUserRequest) (map[string]interface{}, error) {
	requestURL, _ := url.JoinPath(l.BaseURL, "/api/users/", userId)

	var payload bytes.Buffer
	if err := json.NewEncoder(&payload).Encode(data); err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("PATCH", requestURL, &payload)
	logtoToken, err := l.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+logtoToken)
	req.Header.Add("Content-Type", "application/json")

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

func (l LogtoService) FetchUserByToken(token string) (*FetchLogtoUsersResponse, error) {
	jwksURL, err := url.JoinPath(viper.GetString("logto.endpoint"), "/oidc/jwks")
	if err != nil {
		return nil, err
	}

	invalidTokenError := util.
		MakeServiceError(http.StatusUnprocessableEntity).
		AddDetailError("member", "logto token", "invalid logto token")
	_, claims, error := util.ParseTokenWithJWKS(jwksURL, token)
	if error != nil {
		return nil, invalidTokenError.SetMessage("Invalid token " + error.Error())
	}
	// check issuer
	expectedIssuer, _ := url.JoinPath(viper.GetString("logto.endpoint"), "/oidc")
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

	user, err := l.FetchUserById(userId)
	if err != nil {
		return nil, invalidTokenError.SetMessage("Invalid token, error fetching user: " + err.Error())
	}
	return user, nil
}

func MakeLogtoService(endpoint string) LogtoService {
	return LogtoService{
		BaseURL: endpoint,
	}
}

type LogtoUserRole struct {
	TenantId    string `json:"tenantId"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsDefault   bool   `json:"isDefault"`
}
type FetchUserRoleResponse []LogtoUserRole

func (l LogtoService) FetchUserRole(userId string) (FetchUserRoleResponse, error) {
	userRoleURL, err := url.JoinPath(viper.GetString("logto.endpoint"), "/api/users/", userId, "/roles")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("GET", userRoleURL, nil)
	logtoToken, err := l.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+logtoToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf(string(rawBody))
	}

	var body FetchUserRoleResponse
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return nil, err
	}

	return body, nil
}

type FetchUserInfoResponse struct {
	Sub           string                 `json:"sub"`
	Name          string                 `json:"name"`
	Picture       string                 `json:"picture"`
	UpdatedAt     int64                  `json:"updated_at"`
	Username      string                 `json:"username"`
	CreatedAt     int64                  `json:"created_at"`
	Email         string                 `json:"email"`
	EmailVerified bool                   `json:"email_verified"`
	CustomData    map[string]interface{} `json:"custom_data"`
	Roles         []string               `json:"roles"`
}

func (l LogtoService) FetchUserInfo(accessToken string) (FetchUserInfoResponse, error) {
	userInfoEndpointURL, err := url.JoinPath(viper.GetString("logto.endpoint"), "/oidc/me")
	if err != nil {
		return FetchUserInfoResponse{}, err
	}
	req, _ := http.NewRequest("GET", userInfoEndpointURL, nil)

	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FetchUserInfoResponse{}, err
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return FetchUserInfoResponse{}, err
	}

	if res.Status != "200 OK" {
		return FetchUserInfoResponse{}, fmt.Errorf(string(rawBody))
	}

	var body FetchUserInfoResponse
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return FetchUserInfoResponse{}, err
	}

	return body, nil

}

var LogtoServiceApp LogtoService

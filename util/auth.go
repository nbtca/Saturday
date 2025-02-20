package util

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nbtca/saturday/model"
)

// TODO gen key
var key = []byte(genKey())

type Payload struct {
	Who    string
	Member model.Member
	Role   string
}

type Claims struct {
	jwt.StandardClaims
	Payload
}

func CreateToken(payload Payload) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "saturday",
			Subject:   "token",
		},
		Payload: payload,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	return "Bearer " + tokenString, err
}

func CreateBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetTokenString(token string) (string, error) {
	prefix := "Bearer "
	if !strings.HasPrefix(token, prefix) && !strings.HasPrefix(token, strings.ToLower(prefix)) {
		return "", fmt.Errorf("unexpected token: %v", token)
	}
	tokenString := token[len(prefix):]
	return tokenString, nil
}

// parse a jwt token, which should begin with "`Bearer `"
func ParseToken(token string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	tokenString, err := GetTokenString(token)
	if err != nil {
		return nil, nil, err
	}
	t, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})
	return t, Claims, err
}

func ParseTokenWithJWKS(jwksURL string, token string) (*jwt.Token, *jwt.RegisteredClaims, error) {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{}) // See recommended options in the examples directory.
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create JWKS from the given URL.\nError: %s", err)
	}
	tokenString, err := GetTokenString(token)
	if err != nil {
		return nil, nil, err
	}
	Claims := &jwt.RegisteredClaims{}
	t, err := jwt.ParseWithClaims(tokenString, Claims, jwks.Keyfunc)
	return t, Claims, err
}

/*
this is used for testing
"INVALID" to gen invalid token
"EXPIRED" to gen expired token
"NONE" return empty token
*/
func GenToken(auth string, id ...string) string {
	if auth == "INVALID" {
		return "Invalid"
	} else if auth == "EXPIRED" {
		return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTM0NDMxNjIsImRhdGEiOnsidWlkIjoiNmYxZjk3MDItNjZkNi00NDdiLThlNTUtNWYwYzY0N2M4ZDNhIiwicm9sZSI6InVzZXIifSwiaWF0IjoxNjUzMzU2NzYyfQ.ocAxJGhw6Xt2vt7bwGcMeRPLOQOmaspznyu9aI7G670"
	} else if auth == "NONE" {
		return ""
	}
	token, _ := CreateToken(Payload{Who: id[0], Role: auth})
	return token
}

func genKey() string {
	var bytes []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, 24)
	for i := 0; i < 24; i++ {
		result[i] = bytes[rand.Int31()%62]
	}
	return string(result)
}

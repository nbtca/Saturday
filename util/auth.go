package util

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO gen key
var key = []byte("qwejlkqwjelkqwlkqwejlqjelk")

type Payload struct {
	Who  string
	Role string
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
			Subject:   "token",
		},
		Payload: payload,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})
	return token, Claims, err
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
		return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTM0NDMxNjIsImRhdGEiOnsidWlkIjoiNmYxZjk3MDItNjZkNi00NDdiLThlNTUtNWYwYzY0N2M4ZDNhIiwicm9sZSI6InVzZXIifSwiaWF0IjoxNjUzMzU2NzYyfQ.ocAxJGhw6Xt2vt7bwGcMeRPLOQOmaspznyu9aI7G670"
	} else if auth == "NONE" {
		return ""
	}
	token, _ := CreateToken(Payload{Who: id[0], Role: auth})
	return token
}

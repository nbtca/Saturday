package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)
var key = []byte("qwejlkqwjelkqwlkqwejlqjelk")

type Payload struct {
	Id   string
	Role string
}

type Claims struct {
	Id   string
	Role string
	jwt.StandardClaims
}

func CreateToken(payload Payload) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Id:   payload.Id,
		Role: payload.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Subject:   "user token",
		},
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

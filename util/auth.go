package util

import (
	"log"
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
	log.Println(claims.Id)
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

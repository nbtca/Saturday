package middleware

import (
	"net/http"
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

type Role string

const (
	member Role = "member"
	admin  Role = "admin"
)

func TokenInvalid(c *gin.Context){
	c.AbortWithStatusJSON(util.
		MakeServiceError(http.StatusUnprocessableEntity).
		SetMessage("Token Invalid").
		Build())
}

func Auth(role ...Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			TokenInvalid(c)
			return
		}
		token, claims, err := util.ParseToken(tokenString)
		if err != nil || !token.Valid {
			TokenInvalid(c)
			return
		}
		for _, roleObj := range role {
			if string(roleObj) == claims.Role {
				return
			}
		}
		TokenInvalid(c)
	}
}

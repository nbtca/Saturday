package middleware

import (
	"net/http"
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

type Role string

const (
	Member Role = "member"
	Admin  Role = "admin"
)

func Auth(role ...Role) func(c *gin.Context) {
	TokenInvalidErr := util.MakeServiceError(http.StatusUnprocessableEntity).
		SetMessage("Token Invalid")
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(TokenInvalidErr.Build())
			return
		}
		token, claims, err := util.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(TokenInvalidErr.Build())
			return
		}
		for _, roleObj := range role {
			if string(roleObj) == claims.Role {
				return
			}
		}
		c.AbortWithStatusJSON(TokenInvalidErr.Build())
	}
}

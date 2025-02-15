package middleware

import (
	"net/http"
	"slices"

	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type Role string

const (
	Member Role = "member"
	Admin  Role = "admin"
)

func Auth(role ...Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
				SetMessage("not authorized").
				Build())
			return
		}
		// strip bearer
		token = token[7:]
		userinfo, err := service.LogtoServiceApp.FetchUserInfo(token)
		if err != nil {
			c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
				SetMessage("not authorized").
				Build())
			return
		}
		var role string
		if slices.Contains(userinfo.Roles, "Repair Admin") {
			role = "admin"
		} else if slices.Contains(userinfo.Roles, "Repair Member") {
			role = "member"
		}
		if role != "" {
			member, err := service.MemberServiceApp.GetMemberByLogtoId(userinfo.Sub)
			if err != nil {
				c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
					SetMessage("not authorized").
					Build())
			}
			c.Set("id", member.MemberId)
			c.Set("member", member)
			c.Set("role", role)
			c.Set("user", userinfo)
			return
		}
	}
}

// admin is also member
func StepDown(role string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("role", role)
	}
}

package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type Role string

const (
	Member Role = "member"
	Admin  Role = "admin"
)

type AuthContextUser struct {
	UserInfo service.FetchUserInfoResponse
	Role     []string
}

func Auth(acceptableRoles ...Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
				SetMessage("not authorized").
				Build())
			return
		}
		// handel legacy jwt token
		// this is currently used by wechat mini app
		if len(strings.Split(token, ".")) > 1 {
			tokenParsed, claims, err := util.ParseToken(token)
			if err != nil || !tokenParsed.Valid {
				c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
					SetMessage("not authorized").
					Build())
				return
			}
			for _, roleObj := range acceptableRoles {
				if string(roleObj) == claims.Role {
					c.Set("id", claims.Who)
					c.Set("member", claims.Member)
					c.Set("role", claims.Role)
					return
				}
			}
			c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
				SetMessage("not authorized").
				Build())
			return
		}

		// strip bearer
		token, err := util.GetTokenString(token)
		if err != nil {
			c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
				SetMessage("invalid token type").
				Build())
			return
		}
		userinfo, err := service.LogtoServiceApp.FetchUserInfo(token)
		if err != nil {
			c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
				SetMessage("not authorized").
				Build())
			return
		}
		// TODO this is used for backward compatibility, will only use userRoles in the future
		var role string
		userRoles := []string{"client"}
		if slices.Contains(userinfo.Roles, "Repair Admin") {
			userRoles = append(userRoles, "admin")
			role = "admin"
		}
		if slices.Contains(userinfo.Roles, "Repair Member") {
			userRoles = append(userRoles, "member")
			if role == "" {
				role = "member"
			}
		}
		for _, r := range acceptableRoles {
			if slices.Contains(userRoles, string(r)) {
				member, err := service.MemberServiceApp.GetMemberByLogtoId(userinfo.Sub)
				if err != nil {
					c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
						SetMessage("not authorized").
						Build())
				}
				user := AuthContextUser{
					Role:     userRoles,
					UserInfo: userinfo,
				}
				c.Set("id", member.MemberId)
				c.Set("member", member)
				c.Set("role", role)
				c.Set("user", user)
				return
			}
		}
		c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
			SetMessage("not authorized").
			Build())
	}
}

// admin is also member
func StepDown(role string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("role", role)
	}
}

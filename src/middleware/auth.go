package middleware

import (
	"github.com/gin-gonic/gin"
)

type Role string

const (
	member Role = "member"
	admin  Role = "admin"
)

func Auth(role ...Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("id", "2333333333")
		c.Set("role", member)
	}
}

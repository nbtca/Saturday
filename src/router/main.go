package router

import (
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	Router := gin.Default()
	Router.Use(util.ErrorHandler)

	Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	RouterGroup := Router.Group("/")

	MemberGroup := RouterGroup.Group("/members")
	{

		MemberGroup.GET("/", MemberRouterApp.GetByPage)
		MemberGroup.GET("/:MemberId", MemberRouterApp.GetMemberById)

		MemberGroup.POST("/:MemberId", MemberRouterApp.Create)

		MemberGroup.POST("/token", MemberRouterApp.CreateToken)

	}

	return Router
}

package router

import (
	"saturday/src/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	Router := gin.Default()

	Router.Use(middleware.ErrorHandler)

	Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	PublicGroup := Router.Group("/")
	{
		PublicGroup.GET("members/:MemberId", MemberRouterApp.GetPublicMemberById)
		PublicGroup.GET("members/", MemberRouterApp.GetPublicMemberByPage)
		PublicGroup.POST("members/:MemberId/token", MemberRouterApp.CreateToken)
	}

	Router.PUT("member/activate",
		middleware.Auth("member_inactive,admin_inactive"),
		MemberRouterApp.Activate)

	MemberGroup := Router.Group("/")
	MemberGroup.Use(middleware.Auth("member", "admin"))
	{
		MemberGroup.GET("/member", MemberRouterApp.GetMemberById)
		MemberGroup.PUT("/member", MemberRouterApp.Update)
		MemberGroup.PUT("/member/avatar", MemberRouterApp.UpdateAvatar)
	}
	AdminGroup := Router.Group("/")
	AdminGroup.Use(middleware.Auth("admin"))
	{
		AdminGroup.POST("/members/", MemberRouterApp.CreateMany)
		AdminGroup.POST("/members/:MemberId", MemberRouterApp.Create)
		AdminGroup.PATCH("/members/:MemberId", MemberRouterApp.UpdateBasic)
	}

	return Router
}

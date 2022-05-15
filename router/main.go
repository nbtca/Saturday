package router

import (
	"saturday/middleware"

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

		PublicGroup.GET("events/:EventId", EventRouterApp.GetPublicEventById)
		PublicGroup.GET("events/", EventRouterApp.GetPublicEventByPage)

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

		// TODO: set auth requirements
		// allow current member and current user
		MemberGroup.GET("member/events/:EventId", EventRouterApp.GetEventById)
		MemberGroup.GET("member/events/", EventRouterApp.GetEventByPage)

		MemberGroup.POST("member/events/:EventId/accept", EventRouterApp.Accept)
		MemberGroup.DELETE("member/events/:EventId/accept", EventRouterApp.Drop)

		MemberGroup.POST("member/events/:EventId/commit", EventRouterApp.Commit)
		MemberGroup.PATCH("member/events/:EventId/commit", EventRouterApp.AlterCommit)

	}

	AdminGroup := Router.Group("/")
	AdminGroup.Use(middleware.Auth("admin"))
	{
		AdminGroup.POST("/members/", MemberRouterApp.CreateMany)
		AdminGroup.POST("/members/:MemberId", MemberRouterApp.Create)
		AdminGroup.PATCH("/members/:MemberId", MemberRouterApp.UpdateBasic)

		AdminGroup.DELETE("/events/:EventId/commit", EventRouterApp.RejectCommit)
		AdminGroup.POST("/events/:EventId/close", EventRouterApp.Close)

	}

	return Router
}

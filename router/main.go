package router

import (
	"github.com/nbtca/saturday/middleware"

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
		PublicGroup.GET("members", MemberRouterApp.GetPublicMemberByPage)
		PublicGroup.POST("members/:MemberId/token", MemberRouterApp.CreateToken)
		PublicGroup.PATCH("members/:MemberId/logto_id", MemberRouterApp.BindMemberLogtoId)
		PublicGroup.GET("member/token/logto", MemberRouterApp.CreateTokenViaLogtoToken)

		PublicGroup.POST("clients/token/wechat", ClientRouterApp.CreateTokenViaWeChat)

		PublicGroup.GET("events/:EventId", EventRouterApp.GetPublicEventById)
		PublicGroup.GET("events", EventRouterApp.GetPublicEventByPage)

		PublicGroup.GET("setting", SettingRouterApp.GetMiniAppSetting)
	}

	Router.PATCH("member/activate",
		middleware.Auth("member_inactive", "admin_inactive"),
		MemberRouterApp.Activate)

	MemberGroup := Router.Group("/")
	MemberGroup.Use(middleware.Auth("member", "admin"), middleware.StepDown("member"))
	{
		MemberGroup.GET("/member", MemberRouterApp.GetMemberById)
		MemberGroup.PUT("/member", MemberRouterApp.Update)
		MemberGroup.PATCH("/member/avatar", MemberRouterApp.UpdateAvatar)

		MemberGroup.GET("member/events", EventRouterApp.GetMemberEventByPage)
		MemberGroup.GET("member/events/:EventId", EventRouterApp.GetEventById)
		/*
			!!! IMPORTANT !!!
			this middleware is REQUIRED before all handlers that uses event action (except create)
			or there will be panic
		*/
		MemberGroup.Use(middleware.EventActionPreProcess)
		MemberGroup.POST("member/events/:EventId/accept", EventRouterApp.Accept)
		MemberGroup.DELETE("member/events/:EventId/accept", EventRouterApp.Drop)
		MemberGroup.POST("member/events/:EventId/commit", EventRouterApp.Commit)
		MemberGroup.PATCH("member/events/:EventId/commit", EventRouterApp.AlterCommit)

		// MemberGroup.GET("client/:ClientId/events", EventRouterApp.GetEventByClientAndPage)

	}

	AdminGroup := Router.Group("/")
	AdminGroup.Use(middleware.Auth("admin"))
	{
		AdminGroup.POST("/members", MemberRouterApp.CreateMany)
		AdminGroup.POST("/members/:MemberId", MemberRouterApp.Create)
		AdminGroup.PATCH("/members/:MemberId", MemberRouterApp.UpdateBasic)

		AdminGroup.Use(middleware.EventActionPreProcess)
		AdminGroup.DELETE("/events/:EventId/commit", EventRouterApp.RejectCommit)
		AdminGroup.POST("/events/:EventId/close", EventRouterApp.Close)

	}

	ClientGroup := Router.Group("/")
	ClientGroup.Use(middleware.Auth("client"))
	{
		ClientGroup.GET("/client/events/:EventId", EventRouterApp.GetEventById)
		ClientGroup.GET("/client/events", EventRouterApp.GetClientEventByPage)
		ClientGroup.POST("/client/event", EventRouterApp.Create)

		ClientGroup.Use(middleware.EventActionPreProcess)
		ClientGroup.PATCH("/client/events/:EventId", EventRouterApp.Update)
		ClientGroup.DELETE("/client/events/:EventId", EventRouterApp.Cancel)
	}

	Router.POST("/upload", middleware.Auth("member", "admin", "client"), CommonRouterApp.Upload)

	return Router
}

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

		// TODO Restful
		PublicGroup.POST("clients/token/wechat", ClientRouterApp.CreateTokenViaWeChat)

		PublicGroup.GET("events/:EventId", EventRouterApp.GetPublicEventById)
		PublicGroup.GET("events/", EventRouterApp.GetPublicEventByPage)

		PublicGroup.GET("setting", SettingRouterApp.GetMiniAppSetting)
	}

	Router.PUT("member/activate",
		middleware.Auth("member_inactive,admin_inactive"),
		MemberRouterApp.Activate)

	MemberGroup := Router.Group("/")
	MemberGroup.Use(middleware.Auth("member", "admin"), middleware.StepDown("member"))
	{
		MemberGroup.GET("/member", MemberRouterApp.GetMemberById)
		MemberGroup.PUT("/member", MemberRouterApp.Update)
		MemberGroup.PUT("/member/avatar", MemberRouterApp.UpdateAvatar)

		MemberGroup.GET("member/events/", EventRouterApp.GetEventByPage)

		MemberGroup.Use(middleware.EventActionPerProcess)
		MemberGroup.GET("member/events/:EventId", EventRouterApp.GetEventById)
		MemberGroup.POST("member/events/:EventId/accept", EventRouterApp.Accept)
		MemberGroup.DELETE("member/events/:EventId/accept", EventRouterApp.Drop)
		MemberGroup.POST("member/events/:EventId/commit", EventRouterApp.Commit)
		MemberGroup.PATCH("member/events/:EventId/commit", EventRouterApp.AlterCommit)

		MemberGroup.GET("client/:ClientId/events/", EventRouterApp.GetEventByClientAndPage)

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

	ClientGroup := Router.Group("/")
	ClientGroup.Use(middleware.Auth("client"))
	{
		ClientGroup.GET("/client/events/:EventId", EventRouterApp.GetEventById)
		ClientGroup.GET("/client/events/", EventRouterApp.GetClientEventByPage)
		ClientGroup.POST("/client/events", EventRouterApp.Create)
		ClientGroup.PATCH("/client/events/:EventId", EventRouterApp.Update)
		ClientGroup.DELETE("/client/events/:EventId", EventRouterApp.Cancel)
	}

	return Router
}

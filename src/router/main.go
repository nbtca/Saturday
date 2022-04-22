package router

import (
	"gin-example/src/util"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	Router := gin.Default()
	Router.Use(util.ErrorHandler)

	RouterGroup := Router.Group("/")

	MemberGroup := RouterGroup.Group("/members")
	{

		MemberGroup.GET("/", MemberRouterApp.GetByPage)
		MemberGroup.GET("/:MemberId", MemberRouterApp.GetMemberById)

		MemberGroup.POST("/:Member", MemberRouterApp.Create)

		MemberGroup.POST("/token", MemberRouterApp.CreateToken)

	}

	Router.Run()
}

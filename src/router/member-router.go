package router

import (
	"gin-example/src/controller"

	"github.com/gin-gonic/gin"
)

func InitMemerRouter(group *gin.RouterGroup) {
	MemberGroup := group.Group("/member")
	{
		MemberGroup.GET("/:MemberId", controller.MemberControllerApp.GetMemberById)
		MemberGroup.GET("/", controller.MemberControllerApp.GetByPage)
	}

}

package router

import (
	"gin-example/util"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	Router := gin.Default()
	Router.Use(util.ErrorHandler)
	RouterGroup := Router.Group("/")
	{
		InitMemerRouter(RouterGroup)
	}
	Router.Run()
}

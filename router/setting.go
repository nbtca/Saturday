package router

import "github.com/gin-gonic/gin"

type SettingRouter struct{}

func (SettingRouter) GetMiniAppSetting(c *gin.Context) {
	//TODO not implemented
}

var SettingRouterApp = SettingRouter{}

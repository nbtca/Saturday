package router

import "github.com/gin-gonic/gin"

type ClientRouter struct{}

func (ClientRouter) CreateTokenViaWeChat(c *gin.Context) {
	//TODO not implemented
}

var ClientRouterApp = ClientRouter{}

package router

import (
	"saturday/src/model/dto"
	"saturday/src/service"
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

type EventRouter struct{}

func (EventRouter) GetPublicEventById(c *gin.Context) {
	eventId := &dto.EventID{}
	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
		return
	}
	member, err := service.EventServiceApp.GetPublicEventById(eventId.EventID)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (EventRouter) GetEventById(c *gin.Context) {
	eventId := &dto.EventID{}
	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
		return
	}
	member, err := service.EventServiceApp.GetEventById(eventId.EventID)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

var EventRouterApp = EventRouter{}

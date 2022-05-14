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
	event, err := service.EventServiceApp.GetPublicEventById(eventId.EventID)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) GetPublicEventByPage(c *gin.Context) {
	// not implemented
}

func (EventRouter) GetEventById(c *gin.Context) {
	eventId := &dto.EventID{}
	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
		return
	}
	event, err := service.EventServiceApp.GetEventById(eventId.EventID)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) GetEventByPage(c *gin.Context) {
	// not implemented
}

func (EventRouter) Accept(c *gin.Context) {
	eventId := &dto.EventID{}
	memberId := "2333333333"
	// memberId := c.GetString("id")
	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
		return
	}

	event, err := service.EventServiceApp.GetEventById(eventId.EventID)
	if util.CheckError(c, err) {
		return
	}
	if err = service.EventServiceApp.Accept(&event, memberId); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Drop(c *gin.Context) {
	// not implemented
}

func (EventRouter) Commit(c *gin.Context) {
	// not implemented
}

func (EventRouter) AlterCommit(c *gin.Context) {
	// not implemented
}

func (EventRouter) RejectCommit(c *gin.Context) {
	// not implemented
}

func (EventRouter) Close(c *gin.Context) {
	// not implemented
}

var EventRouterApp = EventRouter{}

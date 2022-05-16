package router

import (
	"saturday/model/dto"
	"saturday/service"
	"saturday/util"

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
	//TODO not implemented
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
	//TODO not implemented
}

func (EventRouter) Accept(c *gin.Context) {
	eventId := &dto.EventID{}
	memberId := c.GetString("id")
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
	eventId := &dto.EventID{}
	memberId := c.GetString("id")
	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
		return
	}

	event, err := service.EventServiceApp.GetEventById(eventId.EventID)
	if util.CheckError(c, err) {
		return
	}
	if err = service.EventServiceApp.Drop(&event, memberId); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Commit(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) AlterCommit(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) RejectCommit(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) Close(c *gin.Context) {
	//TODO not implemented
}

var EventRouterApp = EventRouter{}

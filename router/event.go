package router

import (
	"saturday/model"
	"saturday/model/dto"
	"saturday/service"
	"saturday/util"
	action "saturday/util/event-action"

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
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Act(&event, identity, action.Accept); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Drop(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Act(&event, identity, action.Drop); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Commit(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	req := &dto.CommitReq{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Act(&event, identity, action.Commit, req.Content); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
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

func (EventRouter) GetClientEventByPage(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) Create(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) Update(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) Cancel(c *gin.Context) {
	//TODO not implemented
}

func (EventRouter) GetEventByClientAndPage(c *gin.Context) {
	// TODO not implemented
}

var EventRouterApp = EventRouter{}

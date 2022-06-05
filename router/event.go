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
	if err := action.PerformEventAction(&event, identity, action.Accept); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Drop(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	if err := action.PerformEventAction(&event, identity, action.Drop); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Commit(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	req := &dto.CommitReq{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	if err := action.PerformEventAction(&event, identity, action.Commit, req.Content); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) AlterCommit(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	req := &dto.AlterCommitReq{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	if err := action.PerformEventAction(&event, identity, action.Commit, req.Content); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) RejectCommit(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	if err := action.PerformEventAction(&event, identity, action.Reject); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Close(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	if err := action.PerformEventAction(&event, identity, action.Close); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) GetClientEventByPage(c *gin.Context) {

}

func (EventRouter) Create(c *gin.Context) {

}

func (EventRouter) Update(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	req := &dto.UpdateReq{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	if req.Phone != "" {
		event.Phone = req.Phone
	}
	if req.QQ != "" {
		event.QQ = req.QQ
	}
	if req.Phone != "" {
		event.Problem = req.Problem
	}
	if err := action.PerformEventAction(&event, identity, action.Update); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Cancel(c *gin.Context) {
	//TODO not implemented
}

var EventRouterApp = EventRouter{}

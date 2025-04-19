package router

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type EventRouter struct {
	huma huma.API
}

func (EventRouter) GetPublicEventById(c context.Context, input *struct {
	EventID int64 `path:"EventId"`
}) (*util.CommonResponse[model.PublicEvent], error) {
	event, err := service.EventServiceApp.GetPublicEventById(input.EventID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) GetPublicEventByPage(c context.Context, input *struct {
	dto.PageRequest
	Status string `query:"status"`
	Order  string `query:"order" default:"ASC"`
}) (*util.CommonResponse[[]model.PublicEvent], error) {
	events, err := service.EventServiceApp.GetPublicEvents(repo.EventFilter{
		Offset: input.Offset,
		Limit:  input.Limit,
		Status: input.Status,
		Order:  input.Order,
	})
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(events), nil
}

func (er EventRouter) GetEventById(c *gin.Context) {
	eventId := &dto.EventID{}
	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
		return
	}
	event, err := service.EventServiceApp.GetEventById(eventId.EventID)
	identity := util.GetIdentity(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	clientId, err := middleware.GetClientFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if event.MemberId != identity.Member.MemberId && event.ClientId != clientId {
		c.AbortWithStatusJSON(util.MakeServiceError(http.StatusUnauthorized).
			SetMessage("not authorized").
			Build())
		return
	}
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

// return events that is accepted by current member
func (EventRouter) GetMemberEventByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c) // TODO use validator
	if err != nil {
		c.Error(err)
		return
	}
	identity := util.GetIdentity(c)
	status := c.DefaultQuery("status", "")
	order := c.DefaultQuery("order", "ASC")
	events, err := service.EventServiceApp.GetMemberEvents(repo.EventFilter{
		Offset: offset,
		Limit:  limit,
		Status: status,
		Order:  order,
	}, identity.Id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, events)
}

func (EventRouter) Accept(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Accept(&event, identity); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Drop(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Act(&event, identity, util.Drop); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Commit(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	req := &dto.CommitRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	if err := service.EventServiceApp.Act(&event, identity, util.Commit, req.Content); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) AlterCommit(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	req := &dto.AlterCommitRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	if err := service.EventServiceApp.Act(&event, identity, util.AlterCommit, req.Content); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) RejectCommit(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Act(&event, identity, util.Reject); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) Close(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	if err := service.EventServiceApp.Act(&event, identity, util.Close); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (er EventRouter) GetClientEventByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c) // TODO use validator
	if err != nil {
		c.Error(err)
		return
	}
	clientId, err := middleware.GetClientFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	status := c.DefaultQuery("status", "")
	order := c.DefaultQuery("order", "ASC")
	events, err := service.EventServiceApp.GetClientEvents(repo.EventFilter{
		Offset: offset,
		Limit:  limit,
		Status: status,
		Order:  order,
	}, clientId)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, events)
}

func (EventRouter) Create(c *gin.Context) {
	req := &dto.CreateEventRequest{}

	clientId, err := middleware.GetClientFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}
	req.ClientId = clientId

	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	event := &model.Event{
		ClientId:          req.ClientId,
		Model:             req.Model,
		Phone:             req.Phone,
		QQ:                req.QQ,
		ContactPreference: req.ContactPreference,
		Problem:           req.Problem,
	}
	err = service.EventServiceApp.CreateEvent(event)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (er EventRouter) Update(c *gin.Context) {
	rawEvent, _ := c.Get("event")
	event := rawEvent.(model.Event)
	identity := util.GetIdentity(c)
	clientId, err := middleware.GetClientFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}
	identity.ClientId = clientId

	req := &dto.UpdateRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	if req.Phone != "" {
		event.Phone = req.Phone
	}
	if req.QQ != "" {
		event.QQ = req.QQ
	}
	if req.Problem != "" {
		event.Problem = req.Problem
	}
	if req.Model != "" {
		event.Model = req.Model
	}
	if req.ContactPreference != "" {
		event.ContactPreference = req.ContactPreference
	}
	if err := service.EventServiceApp.Act(&event, identity, util.Update); util.CheckError(c, err) {
		return
	}

	c.JSON(200, event)
}

func (er EventRouter) Cancel(c *gin.Context) {
	raw_Event, _ := c.Get("event")
	event := raw_Event.(model.Event)
	identity := util.GetIdentity(c)
	clientId, err := middleware.GetClientFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}
	identity.ClientId = clientId

	if err := service.EventServiceApp.Act(&event, identity, util.Cancel); util.CheckError(c, err) {
		return
	}
	c.JSON(200, event)
}

func (EventRouter) GetEventByClientAndPage(c *gin.Context) {
	// TODO not implemented
}

var EventRouterApp = EventRouter{}

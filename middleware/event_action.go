package middleware

import (
	"saturday/service"
	"saturday/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
 Get event and put to context, and set role to member
 if event's member field equals to id. You are supposed
 to call this before any route that performs event action.
*/
func EventActionPreProcess(c *gin.Context) {
	rawEventId := c.Param("EventId")
	eventId, err := strconv.ParseInt(rawEventId, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(util.MakeValidationError("EventId", nil).
			AddDetailError("Event", "EventId", "Invalid EventId").
			Build())
	}
	role := c.GetString("role")
	id := c.GetString("id")
	event, err := service.EventServiceApp.GetEventById(eventId)
	if util.CheckError(c, err) {
		return
	}
	if role == "member" && event.MemberId == id {
		// set role to current member
		c.Set("role", "member_current")
	}
	c.Set("event", event)
}

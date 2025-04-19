package middleware

import (
	"strconv"

	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

func GetClientFromContext(c *gin.Context) (clientId int64, err error) {
	rawUser := c.Value("user")
	if rawUser == nil {
		return strconv.ParseInt(util.GetIdentity(c).Id, 10, 64)
	}

	user := rawUser.(AuthContextUser)
	if user.UserInfo.Sub != "" {
		logtoId := user.UserInfo.Sub
		client, err := service.ClientServiceApp.CreateClientByLogtoIdIfNotExists(logtoId)
		if err != nil {
			return clientId, err
		}
		util.Logger.Debugf("using logtoId %v for client %v", logtoId, client.ClientId)
		clientId = client.ClientId
	} else {
		clientId, _ = strconv.ParseInt(util.GetIdentity(c).Id, 10, 64)
	}
	return clientId, nil
}

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
	// clientId, err := GetClientFromContext(c)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	// 	return
	// }
	// if role == "client" && event.ClientId == clientId {
	// 	// set role to current client
	// 	c.Set("role", "client_current")
	// }
	c.Set("event", event)
}

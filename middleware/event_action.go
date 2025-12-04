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
Get event and put to context.
This middleware should be added to any route that performs event action.
*/
func EventActionPreProcess(c *gin.Context) {
	rawEventId := c.Param("EventId")
	eventId, err := strconv.ParseInt(rawEventId, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(util.MakeValidationError("EventId", nil).
			AddDetailError("Event", "EventId", "Invalid EventId").
			Build())
	}
	event, err := service.EventServiceApp.GetEventById(eventId)
	if util.CheckError(c, err) {
		return
	}
	c.Set("event", event)
}

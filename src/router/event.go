package router

type EventRouter struct{}

// func (EventRouter) GetPublicEventById(c *gin.Context) {
// 	eventId := &dto.EventID{}
// 	if err := util.BindAll(c, eventId); util.CheckError(c, err) {
// 		return
// 	}
// 	member, err := service.EventService.GetPublicEventById(eventId.EventID)
// 	if util.CheckError(c, err) {
// 		return
// 	}
// 	c.JSON(200, member)
// }

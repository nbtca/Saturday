package service

import (
	"log"
	"net/http"
	"saturday/src/model"
	"saturday/src/repo"
	"saturday/src/util"
)

type EventService struct{}

func (service EventService) GetEventById(id int64) (model.Event, error) {
	event, err := repo.GetEventById(id)
	if err != nil || event.EventId == 0 {
		return model.Event{}, util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed")
	}
	return event, nil
}

func (service EventService) GetPublicEventById(id int64) (model.PublicEvent, error) {
	event, err := service.GetEventById(id)
	if err != nil {
		return model.PublicEvent{}, err
	}
	return model.CreatePublicEvent(event), nil
}

func (service EventService) Accept(event model.Event, memberId string) (model.Event, error) {
	// judge whether the event is acceptable
	if event.Status != "open" {
		return model.Event{},
			util.MakeServiceError(http.StatusUnprocessableEntity).
				SetMessage("This event is not open")
	}
	event.Status = "accepted"
	event.MemberId = memberId
	acceptLog := model.EventLog{
		EventId:  event.EventId,
		MemberId: memberId,
		Action:   "Accept",
	}
	_, err := repo.UpdateEvent(event)
	if err != nil {
		// TODO  wrap error
		log.Println(err)
		return model.Event{}, err
	}
	if err = repo.CreateEventLog(&acceptLog); err != nil {
		// TODO  wrap error
		log.Println(err)
		return model.Event{}, err
	}
	event.Logs = append(event.Logs, acceptLog)
	return event, nil
}

var EventServiceApp = EventService{}

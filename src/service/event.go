package service

import (
	"net/http"
	"saturday/src/model"
	"saturday/src/repo"
	"saturday/src/util"
)

type EventService struct{}

func (service EventService) GetEventById(id int64) (model.Event, error) {
	event, err := repo.GetEventById(id)
	if err != nil {
		return model.Event{}, err
	}
	if event.EventId == 0 {
		error := util.MakeServiceError(
			http.StatusUnprocessableEntity).
			SetMessage("Validation Failed")
		return event, error
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

var EventServiceApp = EventService{}

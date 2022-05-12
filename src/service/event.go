package service

import (
	"saturday/src/model"
)

type EventService struct{}

func (service *EventService) GetEventById(id int64) (model.Event, error) {
	// event, err := repo.GetEventById(id)
	// if err != nil {
	// 	return model.Event{}, err
	// }
	// if event == (model.Event{}) {
	// 	error := util.MakeServiceError(
	// 		http.StatusUnprocessableEntity).
	// 		SetMessage("Validation Failed")
	// 	return event, error
	// } else {
	// 	return event, nil
	// }
	return model.Event{}, nil
}

func (service *EventService) GetPublicEventById(id int64) (model.PublicEvent, error) {
	event, err := service.GetEventById(id)
	if err != nil {
		return model.PublicEvent{}, err
	}
	return model.CreatePublicEvent(event), nil
}

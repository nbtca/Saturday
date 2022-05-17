package service

import (
	"net/http"
	"saturday/model"
	"saturday/repo"
	"saturday/util"
)

type EventService struct {
	model.Event
}

// not used
func MakeEventService(id int64) (*EventService, error) {
	event, err := EventServiceApp.GetEventById(id)
	if err != nil {
		return nil, err
	}
	return &EventService{event}, nil
}

func (service EventService) GetEventById(id int64) (model.Event, error) {
	event, err := repo.GetEventById(id)
	if err != nil || event.EventId == 0 {
		util.Logger.Error(err)
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

// func (service EventService) Accept(event *model.Event, identity model.Identity) error {
// 	err := action.PerformEventAction(event, identity, action.Accept)
// 	return err
// }

// func (service EventService) Drop(event *model.Event, identity model.Identity) error {
// 	err := action.PerformEventAction(event, identity, action.Drop)
// 	return err
// }

var EventServiceApp = EventService{}

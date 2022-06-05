package service

import (
	"fmt"
	"net/http"
	"saturday/model"
	"saturday/repo"
	"saturday/util"
	au "saturday/util/event-action"
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
	if err != nil {
		util.Logger.Error(err)
		return model.Event{}, util.MakeInternalServerError()
	}
	if event.EventId == 0 {
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

func (service EventService) GetPublicEvents(offset uint64, limit uint64) ([]model.PublicEvent, error) {
	events, err := repo.GetEvents(offset, limit)
	if err != nil {
		return nil, err
	}
	publicEvents := make([]model.PublicEvent, len(events))
	for i, v := range events {
		publicEvents[i] = model.CreatePublicEvent(v)
	}
	return publicEvents, nil
}

/*
 this function validates the action and then perform action to the event.
 it also persists the event and event log.
*/
func (service EventService) Act(event *model.Event, identity model.Identity, action au.Action, description ...string) error {
	handler := au.EventActionMap[action]
	handler.Init(event, identity)
	for _, d := range description {
		handler.Description = fmt.Sprint(handler.Description, d)
	}
	if err := handler.ValidateAction(); err != nil {
		return err
	}
	return handler.Handle()
}

var EventServiceApp = EventService{}

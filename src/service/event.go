package service

import (
	"log"
	"net/http"
	"saturday/src/model"
	"saturday/src/repo"
	"saturday/src/util"
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

func (service EventService) Accept(event *model.Event, memberId string) error {
	// judge whether the event is acceptable
	if event.Status != "open" {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("This event is not open")
	}
	event.Status = "accepted"
	event.MemberId = memberId
	acceptLog := model.EventLog{
		EventId:  event.EventId,
		MemberId: memberId,
		Action:   "Accept",
	}
	_, err := repo.UpdateEvent(*event)
	if err != nil {
		// TODO  wrap error
		log.Println(err)
		return err
	}
	if err = repo.CreateEventLog(&acceptLog); err != nil {
		// TODO  wrap error
		log.Println(err)
		return err
	}
	event.Logs = append(event.Logs, acceptLog)
	return nil
}

var EventServiceApp = EventService{}

package service

import (
	"fmt"
	"net/http"
	"net/rpc"
	"net/url"
	"os"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"

	"gopkg.in/gomail.v2"
)

type EventService struct{}

func (service EventService) GetEventById(id int64) (model.Event, error) {
	event, err := repo.GetEventById(id)
	if err != nil {
		return model.Event{}, util.MakeInternalServerError()
	}
	if event.EventId == 0 {
		return model.Event{}, util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed")
	}
	return event, nil
}

func (service EventService) GetMemberEvents(f repo.EventFilter, memberId string) ([]model.Event, error) {
	return repo.GetMemberEvents(f, memberId)
}

func (service EventService) GetClientEvents(f repo.EventFilter, clientId string) ([]model.Event, error) {
	return repo.GetClientEvents(f, clientId)
}

func (service EventService) GetPublicEventById(id int64) (model.PublicEvent, error) {
	event, err := service.GetEventById(id)
	if err != nil {
		return model.PublicEvent{}, err
	}
	return model.CreatePublicEvent(event), nil
}

func (service EventService) GetPublicEvents(f repo.EventFilter) ([]model.PublicEvent, error) {
	events, err := repo.GetEvents(f)
	if err != nil {
		return nil, err
	}
	publicEvents := make([]model.PublicEvent, len(events))
	for i, v := range events {
		publicEvents[i] = model.CreatePublicEvent(v)
	}
	return publicEvents, nil
}

func (service EventService) CreateEvent(event *model.Event) error {
	// insert event
	if err := repo.CreateEvent(event); err != nil {
		return err
	}
	identity := model.Identity{
		Id:   fmt.Sprint(event.ClientId),
		Role: "client",
	}
	// insert event status and event log
	if err := service.Act(event, identity, util.Create); err != nil {
		return err
	}
	return nil
}

func (service EventService) SendActionNotify(event *model.Event, subject string) error {
	if event == nil {
		return util.MakeInternalServerError()
	}
	service.SendActionNotifyViaPushDeer(event, subject)
	if err := service.SendActionNotifyViaRPC(&model.EventActionNotifyRequest{
		Subject:   subject,
		Model:     event.Model,
		Problem:   event.Problem,
		Link:      "A Link to Sunday",
		GmtCreate: event.GmtCreate,
	}); err != nil {
		return service.SendActionNotifyViaMail(event, subject)
	}
	return nil
}

// A temporary function to send action notify via PushDeer
func (service EventService) SendActionNotifyViaPushDeer(event *model.Event, subject string) error {
	pushKey := "PDU6809T95RlAWSw5pMFI2qy5EIFGjiVEapfV6Qs"
	pushAPI, _ := url.Parse("https://api2.pushdeer.com/message/push")
	params := url.Values{}
	params.Add("pushkey", pushKey)
	params.Add("text", fmt.Sprintf("型号：%s，问题描述：%s，创建时间：%s", event.Model, event.Problem, event.GmtCreate))
	pushAPI.RawQuery = params.Encode()
	_, err := http.Get(pushAPI.String())
	return err
}

func (service EventService) SendActionNotifyViaRPC(req *model.EventActionNotifyRequest) error {
	address := os.Getenv("RPC_ADDRESS")
	if address == "" {
		return fmt.Errorf("RPC_ADDRESS is not set")
	}
	conn, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return err
	}
	// req := model.EventActionNotifyRequest{
	// 	Subject:   subject,
	// 	Model:     event.Model,
	// 	Problem:   event.Problem,
	// 	Link:      "A Link to Sunday",
	// 	GmtCreate: event.GmtCreate,
	// }
	res := model.EventActionNotifyResponse{}
	if err = conn.Call("Notify.EventActionNotify", req, &res); err != nil {
		util.Logger.Error(err)
		return err
	}
	if !res.Success {
		return fmt.Errorf("failed to send action notify via rpc")
	}
	return nil
}

func (service EventService) SendActionNotifyViaMail(event *model.Event, subject string) error {
	m := gomail.NewMessage()
	receiverAddress := os.Getenv("MAIL_RECEIVER_ADDRESS")
	if receiverAddress == "" {
		return fmt.Errorf("MAIL_RECEIVER_ADDRESS is not set")
	}
	m.SetHeader("To", receiverAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", fmt.Sprintf(
		`<div>
  <span style="padding-right:10px;">型号:</span>
  <span>%s</span>
</div>
<div>
  <span style="padding-right:10px;">问题描述:</span>
  <span>%s</span>
</div>
<div>
  <span style="padding-right:10px;">创建时间:</span>
  <span>%s</span>
</div>
<div style="padding-top:10px;">
  <a href="https://repair.nbtca.space">在 Sunday 中处理</a>
</div>`, event.Model, event.Problem, event.GmtCreate))

	if err := util.SendMail(m); err != nil {
		return util.MakeInternalServerError().SetMessage("fail on mail")
	}
	return nil
}

/*
this function validates the action and then perform action to the event.
it also persists the event and event log.
*/
func (service EventService) Act(event *model.Event, identity model.Identity, action util.Action, description ...string) error {
	handler := util.MakeEventActionHandler(action, event, identity)
	if err := handler.ValidateAction(); err != nil {
		return err
	}
	for _, d := range description {
		handler.Description = fmt.Sprint(handler.Description, d)
	}

	log := handler.Handle()
	// persist event
	if err := repo.UpdateEvent(event, &log); err != nil {
		return err
	}
	// append log
	event.Logs = append(event.Logs, log)
	return nil
}

var EventServiceApp = EventService{}

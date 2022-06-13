package util_test

import (
	"fmt"
	"log"
	"saturday/model"
	"saturday/util"
	"testing"
)

func TestCreateToken(t *testing.T) {
	j, _ := util.CreateToken(util.Payload{Who: "123", Role: "member"})
	_, claims, _ := util.ParseToken(j)
	if claims.Who != "123" || claims.Role != "member" {
		t.Error("测试失败")
	}
}

func TestReadCsvFile(t *testing.T) {
	records := util.ReadCsvFile("testdata/event-action.csv")
	for _, v := range records {
		fmt.Println(v)
	}
}
func TestGetCsvMap(t *testing.T) {
	records, err := util.GetCsvMap("testdata/event-action.csv")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range records {
		fmt.Println(v)
	}
}

func TestEventAction(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/event-action.csv")
	if err != nil {
		t.Error(err)
	}
	// event := model.Event{}
	for _, v := range rawAPITestCase {
		t.Run(v["CaseId"], func(t *testing.T) {
			event := &model.Event{
				Status: v["event.status"],
			}
			actor := model.Identity{
				Id:   "2333333333",
				Role: v["actor.role"],
			}
			handler := util.MakeEventActionHandler(util.Action(v["action"]), event, actor)
			err := handler.ValidateAction()
			if err != nil {
				if v["error"] != "X" {
					se, _ := util.IsServiceError(err)
					t.Errorf("%s: %s", v["CaseId"], se.Body.Message)
				}
			} else {
				if v["error"] == "X" {
					log.Println(event, actor,v["action"])
					t.Errorf("error expected")
				}
				log := handler.Handle()
				statusExpected := v["out_event.status"]
				if event.Status != statusExpected {
					t.Errorf("invalid event.status: expected:%v, got:%v", statusExpected, event.Status)
				}
				action := v["out_event.action"]
				if action != "" && log.Action != action {
					t.Errorf("invalid event.action: expected:%v, got:%v", action, log.Action)
				}
				logMemberId := v["out_event.member_id"]
				if logMemberId == "actor.id" {
					logMemberId = actor.Id
				}
				if logMemberId != "" && log.MemberId != logMemberId {
					t.Errorf("invalid event.member_id: expected:%v, got:%v", logMemberId, log.MemberId)
				}
				memberId := v["out_event.memberId"]
				if memberId == "actor.id" {
					memberId = actor.Id
				}
				if memberId != "" && event.MemberId != memberId {
					t.Errorf("invalid event.memberId: expected:%v, got:%v", memberId, event.MemberId)
				}
				closedBy := v["out_event.closed_by"]
				if closedBy == "actor.id" {
					closedBy = actor.Id
				}
				if closedBy != "" && event.ClosedBy != closedBy {
					t.Errorf("invalid event.closed_by: expected:%v, got:%v", closedBy, event.ClosedBy)
				}
			}
		})
	}
}

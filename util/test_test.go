package util_test

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// func TestCreateToken(t *testing.T) {
// 	j, _ := util.CreateToken(util.Payload{Who: "123", Role: "member"})
// 	_, claims, _ := util.ParseToken(j)
// 	if claims.Who != "123" || claims.Role != "member" {
// 		t.Error("测试失败")
// 	}
// }

// func TestReadCsvFile(t *testing.T) {
// 	records := util.ReadCsvFile("testdata/event-action.csv")
// 	for _, v := range records {
// 		fmt.Println(v)
// 	}
// }
// func TestGetCsvMap(t *testing.T) {
// 	records, err := util.GetCsvMap("testdata/event-action.csv")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	for _, v := range records {
// 		fmt.Println(v)
// 	}
// }

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		util.Logger.Warning("Error loading .env file")
	}
	m.Run()
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
				logMemberId := v["out_event.memberId"]
				if logMemberId == "actor.id" {
					logMemberId = actor.Id
				}
				if logMemberId != "" && log.MemberId != logMemberId {
					t.Errorf("invalid event.memberId: expected:%v, got:%v", logMemberId, log.MemberId)
				}
				memberId := v["out_event.memberId"]
				if memberId == "actor.id" {
					memberId = actor.Id
				}
				if memberId != "" && event.MemberId != memberId {
					t.Errorf("invalid event.memberId: expected:%v, got:%v", memberId, event.MemberId)
				}
				closedBy := v["out_event.closedBy"]
				if closedBy == "actor.id" {
					closedBy = actor.Id
				}
				if closedBy != "" && event.ClosedBy != closedBy {
					t.Errorf("invalid event.closedBy: expected:%v, got:%v", closedBy, event.ClosedBy)
				}
			}
		})
	}
}

var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

func TestParseTokenWithoutBearer(t *testing.T) {
	_, _, err := util.ParseToken(tokenString)
	log.Print(err)
	if err == nil {
		t.Error("no error at missing `Bearer `")
	}
}

func TestParseToken(t *testing.T) {
	token := "Bearer " + tokenString
	_, claims, err := util.ParseToken(token)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(claims)
}

func TestParseTokenWithJWKS(t *testing.T) {
	token := "Bearer eyJhbGciOiJFUzM4NCIsInR5cCI6ImF0K2p3dCIsImtpZCI6Im9VU0hpdWNoNkpGUS1yaGRiTnFvLVRrVy1VRmpudmtSako3aWw1dFdOYU0ifQ.eyJqdGkiOiI4VW10UWVlMjVvZzRlSGc4cl9NUHMiLCJzdWIiOiJjaG16MWl0ejgzcXEiLCJpYXQiOjE2OTg3NTcxMDUsImV4cCI6MTY5ODc2MDcwNSwic2NvcGUiOiIiLCJjbGllbnRfaWQiOiJoMmVqa2tmd2R0ampwZW1iMDIxcm8iLCJpc3MiOiJodHRwczovL2F1dGguYXBwLm5idGNhLnNwYWNlL29pZGMiLCJhdWQiOiJodHRwczovL2FwaS5uYnRjYS5zcGFjZS92MiJ9.uUzXk8zERRhWtWFMnLcLGDF8ZQl-PoSWVWv6MnCjHb1q5P1aHlKVRx2RmSjDr2Nm7n0JZIXsSVQrDXhsB0J64qi2gI4Xuvu3pe11FIpeVxHLY7ObpDzyaeRBHc26P2Lo"
	jwksURL, err := url.JoinPath(viper.GetString("LOGTO_ENDPOINT"), "/oidc/jwks")
	if err != nil {
		t.Error(err)
		return
	}
	_, claims, err := util.ParseTokenWithJWKS(jwksURL, token)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(claims)
}

func TestSendMail(t *testing.T) {

	m := gomail.NewMessage()
	util.InitDialer()
	receiverAddress := viper.GetString("testing.mail.receiver_address")
	if receiverAddress == "" {
		t.Error("receiver_address is not set")
	}
	m.SetHeader("To", receiverAddress)
	m.SetHeader("Subject", "维修状态更新(#12): ")
	m.SetBody("text/html", fmt.Sprintf(
		`
		<h3>新的状态为：accepted</h3>
		<div>
  		<span style="padding-right:10px;">问题描述:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">型号:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">创建时间:</span>
  		<span>%s</span>
		</div>
		<div>
  		<span style="padding-right:10px;">手机:</span>
  		<span>13488888888</span>
		</div>
		<div>
  		<span style="padding-right:10px;">QQ:</span>
  		<span>2359845989</span>
		</div>
		<div style="padding-top:10px;">
  		<a href="http://github.com/nbtca/repair-tickets/issues/%v">在 nbtca/repair-tickets 中处理</a>
		</div>
`, "MacBook", "clean", "2021-10-10", 10))
	// 	m.SetBody("text/html", fmt.Sprintf(
	// 		`<div>
	//   <span style="padding-right:10px;">型号:</span>
	//   <span>%s</span>
	// </div>
	// <div>
	//   <span style="padding-right:10px;">问题描述:</span>
	//   <span>%s</span>
	// </div>
	// <div>
	//   <span style="padding-right:10px;">创建时间:</span>
	//   <span>%s</span>
	// </div>
	// <div style="padding-top:10px;">
	//   <a href="https://repair.nbtca.space">在 Sunday 中处理</a>
	// </div>`, "MacBook", "clean", "2021-10-10"))

	if err := util.SendMail(m); err != nil {
		t.Error(err)
	}
}

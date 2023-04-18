package router_test

import (
	"strconv"
	"testing"

	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

func TestGetEventById(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/get_event_by_id.csv")
	if err != nil {
		t.Error(err)
	}
	for _, rawCase := range rawAPITestCase {
		t.Run(rawCase["CaseId"], func(t *testing.T) {
			code, _ := strconv.Atoi(rawCase["code"])
			// auth := rawCase["Authorization"]
			testCase := APITestCase{
				"success",
				Request{
					"GET",
					"/events/" + rawCase["eventId"],
					// GenToken(auth, "2333333333"),
					"",
					gin.H{},
				},
				Response{
					code,
					gin.H{
						"eventId":  1,
						"clientId": 1,
						"model":    "7590",
						"problem":  "hackintosh",
						"member": gin.H{
							"memberId":    "2333333333",
							"alias":       "滑稽",
							"name":        "滑稽",
							"section":     "计算机233",
							"role":        "member",
							"profile":     "relaxing",
							"avatar":      "",
							"createdBy":   "0000000000",
							"gmtCreate":   IGNORE,
							"gmtModified": IGNORE,
						},
						"closedBy": gin.H{
							"memberId":    "0000000000",
							"alias":       "管理",
							"name":        "管理",
							"section":     "计算机000",
							"role":        "admin",
							"profile":     "",
							"avatar":      "",
							"createdBy":   "",
							"gmtCreate":   IGNORE,
							"gmtModified": IGNORE,
						},
						"status":      "accepted",
						"logs":        IGNORE,
						"gmtCreate":   IGNORE,
						"gmtModified": IGNORE,
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				testCase.Response.Body = gin.H{
					"message": rawCase["error_message"],
				}
			}
			err = testCase.Test()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCreateEvent(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/create_event.csv")
	if err != nil {
		t.Error(err)
	}
	for _, rawCase := range rawAPITestCase {
		t.Run(rawCase["CaseId"], func(t *testing.T) {
			code, _ := strconv.Atoi(rawCase["code"])
			auth := rawCase["Authorization"]
			testCase := APITestCase{
				"success",
				Request{
					"POST",
					"/client/events",
					util.GenToken(auth, "1"),
					gin.H{
						"model":             rawCase["model"],
						"phone":             rawCase["phone"],
						"qq":                rawCase["qq"],
						"contactPreference": rawCase["contactPreference"],
						"problem":           rawCase["problem"],
					},
				},
				Response{
					code,
					gin.H{
						"eventId":           "IGNORE",
						"clientId":          1,
						"model":             rawCase["model"],
						"phone":             rawCase["phone"],
						"qq":                rawCase["qq"],
						"contactPreference": rawCase["contactPreference"],
						"problem":           rawCase["problem"],
						"memberId":          "",
						"closedBy":          "",
						"status":            "open",
						"logs":              "IGNORE",
						"gmtCreate":         "IGNORE",
						"gmtModified":       "IGNORE",
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				message := ""
				if rawCase["error_message"] == "Token_Invalid" {
					message = "Token Invalid"
				}
				if rawCase["error_message"] == "Validation_Failed" {
					message = "Validation Failed"
				}
				testCase.Response.Body = gin.H{
					"message": message,
				}
			}
			err := testCase.Test()
			if err != nil {
				t.Error(err)
			}

		})
	}
}
func BenchmarkGetEventById(b *testing.B) {
	testCase := APITestCase{
		"success",
		Request{
			"GET",
			"/events/1",
			"",
			gin.H{},
		},
		Response{},
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testCase.Run()
		}
	})
}

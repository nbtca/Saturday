package router_test

import (
	"saturday/util"
	"strconv"
	"testing"

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
					"/events/" + rawCase["event_id"],
					// GenToken(auth, "2333333333"),
					"",
					gin.H{},
				},
				Response{
					code,
					gin.H{
						"event_id":  1,
						"client_id": 1,
						"model":     "7590",
						"problem":   "hackintosh",
						"member": gin.H{
							"member_id":    "2333333333",
							"alias":        "滑稽",
							"name":         "滑稽",
							"section":      "计算机233",
							"role":         "member",
							"profile":      "relaxing",
							"phone":        "12356839487",
							"qq":           "123456",
							"avatar":       "",
							"created_by":   "0000000000",
							"gmt_create":   "2022-04-23 15:49:59",
							"gmt_modified": "2022-04-30 17:29:46",
						},
						"closed_by": gin.H{
							"member_id":    "0000000000",
							"alias":        "管理",
							"name":         "管理",
							"section":      "计算机000",
							"role":         "admin",
							"profile":      "",
							"phone":        "",
							"qq":           "",
							"avatar":       "",
							"created_by":   "",
							"gmt_create":   "IGNORE",
							"gmt_modified": "IGNORE",
						},
						"status":       "accepted",
						"logs":         "IGNORE",
						"gmt_create":   "IGNORE",
						"gmt_modified": "IGNORE",
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
						"model":              rawCase["model"],
						"phone":              rawCase["phone"],
						"qq":                 rawCase["qq"],
						"contact_preference": rawCase["contact_preference"],
						"problem":            rawCase["problem"],
					},
				},
				Response{
					code,
					gin.H{
						"event_id":           "IGNORE",
						"client_id":          1,
						"model":              rawCase["model"],
						"phone":              rawCase["phone"],
						"qq":                 rawCase["qq"],
						"contact_preference": rawCase["contact_preference"],
						"problem":            rawCase["problem"],
						"member_id":          "",
						"closed_by":          "",
						"status":             "open",
						"logs":               "IGNORE",
						"gmt_create":         "IGNORE",
						"gmt_modified":       "IGNORE",
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

package router_test

import (
	"saturday/util"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPublicMemberById(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/get_public_member_by_id.csv")
	if err != nil {
		t.Error(err)
	}

	for _, rawCase := range rawAPITestCase {
		t.Run(rawCase["CaseId"], func(t *testing.T) {
			code, _ := strconv.Atoi(rawCase["code"])
			testCase := APITestCase{
				"success",
				Request{
					"GET",
					"/members/" + rawCase["id"],
					"",
					gin.H{},
				},
				Response{
					code,
					gin.H{
						"memberId":    "2333333333",
						"alias":       "滑稽",
						"role":        "member",
						"profile":     "relaxing",
						"avatar":      "",
						"createdBy":   "0000000000",
						"gmtCreate":   "2022-04-23 15:49:59",
						"gmtModified": "2022-04-30 17:29:46",
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				testCase.Response.Body = gin.H{
					"message": "Validation Failed",
				}
			}
			err := testCase.Test()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetMemberById(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/get_member_by_id.csv")
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
					"GET",
					"/member",
					util.GenToken(auth, "2333333333"),
					gin.H{},
				},
				Response{
					code,
					gin.H{
						"memberId":    "2333333333",
						"alias":       "滑稽",
						"name":        "滑稽",
						"section":     "计算机233",
						"role":        "member",
						"profile":     "relaxing",
						"phone":       "12356839487",
						"qq":          "123456",
						"avatar":      "",
						"createdBy":   "0000000000",
						"gmtCreate":   "2022-04-23 15:49:59",
						"gmtModified": "2022-04-30 17:29:46",
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				testCase.Response.Body = gin.H{
					"message": "Token Invalid",
				}
			}
			err = testCase.Test()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUpdateMember(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/update_member.csv")
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
					"PUT",
					"/member",
					util.GenToken(auth, "2333333333"),
					gin.H{
						"alias":    rawCase["alias"],
						"profile":  rawCase["profile"],
						"phone":    rawCase["phone"],
						"qq":       rawCase["qq"],
						"avatar":   rawCase["avatar"],
						"password": rawCase["password"],
					},
				},
				Response{
					code,
					gin.H{
						"memberId":    "2333333333",
						"alias":       rawCase["alias"],
						"name":        "滑稽",
						"section":     "计算机233",
						"role":        "member",
						"profile":     rawCase["profile"],
						"phone":       rawCase["phone"],
						"qq":          rawCase["qq"],
						"avatar":      rawCase["avatar"],
						"createdBy":   "0000000000",
						"gmtCreate":   "2022-04-23 15:49:59",
						"gmtModified": "2022-04-30 17:29:46",
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				testCase.Response.Body = gin.H{
					"message": strings.Replace(rawCase["error_message"], "_", " ", -1),
				}
			}
			err = testCase.Test()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCreateMember(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/create_member.csv")
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
					"/members/" + rawCase["memberId"],
					util.GenToken(auth, "0000000000"),
					gin.H{
						"alias":   rawCase["alias"],
						"name":    rawCase["name"],
						"section": rawCase["section"],
						"role":    rawCase["role"],
						"profile": rawCase["profile"],
						"phone":   rawCase["phone"],
						"qq":      rawCase["qq"],
						"avatar":  rawCase["avatar"],
					},
				},
				Response{
					code,
					gin.H{
						"memberId":    rawCase["memberId"],
						"alias":       rawCase["alias"],
						"name":        rawCase["name"],
						"section":     rawCase["section"],
						"role":        rawCase["role"],
						"profile":     rawCase["profile"],
						"phone":       rawCase["phone"],
						"qq":          rawCase["qq"],
						"avatar":      rawCase["avatar"],
						"createdBy":   "0000000000",
						"gmtCreate":   "2022-04-23 15:49:59",
						"gmtModified": "2022-04-30 17:29:46",
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				testCase.Response.Body = gin.H{
					"message": strings.Replace(rawCase["error_message"], "_", " ", -1),
				}
			}
			err = testCase.Test()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCreateMemberToken(t *testing.T) {
	rawAPITestCase, err := util.GetCsvMap("testdata/create_member_token.csv")
	if err != nil {
		t.Error(err)
	}
	for _, rawCase := range rawAPITestCase {
		t.Run(rawCase["CaseId"], func(t *testing.T) {
			code, _ := strconv.Atoi(rawCase["code"])
			testCase := APITestCase{
				"success",
				Request{
					"POST",
					"/members/" + rawCase["id"] + "/token",
					"",
					gin.H{
						"password": rawCase["password"],
					},
				},
				Response{
					code,
					gin.H{
						"memberId":    "2333333333",
						"alias":       "滑稽",
						"name":        "滑稽",
						"section":     "计算机233",
						"role":        "member",
						"profile":     "relaxing",
						"phone":       "12356839487",
						"qq":          "123456",
						"avatar":      "",
						"createdBy":   "0000000000",
						"gmtCreate":   "IGNORE",
						"gmtModified": "IGNORE",
						"token":       "IGNORE",
					},
				},
			}
			if rawCase["success"] != "TRUE" {
				testCase.Response.Body = gin.H{
					"message": "Validation Failed",
				}
			}
			err := testCase.Test()
			if err != nil {
				t.Error(err)
			}

		})
	}
}

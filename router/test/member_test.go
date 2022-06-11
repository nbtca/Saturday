package router_test

import (
	"saturday/util"
	"strconv"
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
						"member_id":    "2333333333",
						"alias":        "滑稽",
						"role":         "member",
						"profile":      "relaxing",
						"avatar":       "",
						"created_by":   "0000000000",
						"gmt_create":   "2022-04-23 15:49:59",
						"gmt_modified": "2022-04-30 17:29:46",
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
					GenToken(auth, "2333333333"),
					gin.H{},
				},
				Response{
					code,
					gin.H{
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

// func TestMemberActive(t *testing.T) {
// 	for _, data := range MemberActiveData {
// 		t.Run(data.Name, func(t *testing.T) {
// 			err := DataHandler(data)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }

// func TestUpdateMember(t *testing.T) {
// 	for _, data := range UpdateMemberData {
// 		t.Run(data.Name, func(t *testing.T) {
// 			err := DataHandler(data)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }

// func TestUpdateMemberBasic(t *testing.T) {
// 	for _, data := range UpdateMemberBasicInfoData {
// 		t.Run(data.Name, func(t *testing.T) {
// 			err := DataHandler(data)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }
// func TestCreateMember(t *testing.T) {
// 	for _, data := range CreateMemberData {
// 		t.Run(data.Name, func(t *testing.T) {
// 			err := DataHandler(data)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }

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
						"gmt_create":   "IGNORE",
						"gmt_modified": "IGNORE",
						"token":        "IGNORE",
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

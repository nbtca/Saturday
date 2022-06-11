package router_test

import "github.com/gin-gonic/gin"

var GetPublicMemberData = []APITestCase{
	{
		"success",
		Request{
			"GET",
			"/members/2333333333",
			"",
			nil,
		},
		Response{
			200,
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
	},
}

var CreateMemberTokenData = []APITestCase{
	{
		"success",
		Request{
			"POST",
			"/members/2333333333/token",
			"",
			gin.H{
				"password": "123456",
			},
		},
		Response{
			200,
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
				"token":        "not implemented",
			},
		},
	},
}

var GetMemberData = []APITestCase{
	{
		"success",
		Request{
			"GET",
			"/member",
			"",
			nil,
		},
		Response{
			200,
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
	},
}

var MemberActiveData = []APITestCase{
	{
		"success",
		Request{
			"PUT",
			"/member/activate",
			"",
			gin.H{
				"alias":    "滑稽",
				"phone":    "12356839487",
				"qq":       "123456",
				"password": "123456",
			},
		},
		Response{
			200,
			gin.H{
				"member_id":    "2333333333",
				"alias":        "滑da稽",
				"name":         "滑稽",
				"section":      "计算机233",
				"role":         "member",
				"profile":      "want to relax",
				"phone":        "12356839487",
				"qq":           "123456",
				"avatar":       "",
				"created_by":   "0000000000",
				"gmt_create":   "2022-04-23 15:49:59",
				"gmt_modified": "2022-04-30 17:29:46",
			},
		},
	},
}

var UpdateMemberData = []APITestCase{
	{
		"success",
		Request{
			"PUT",
			"/member",
			"",
			gin.H{
				"member_id": "2333333333",
				"alias":     "滑da稽",
				"name":      "滑稽",
				"profile":   "want to relax",
				"phone":     "12356839487",
				"qq":        "123456",
			},
		},
		Response{
			200,
			gin.H{
				"member_id":    "2333333333",
				"alias":        "滑da稽",
				"name":         "滑稽",
				"section":      "计算机233",
				"role":         "member",
				"profile":      "want to relax",
				"phone":        "12356839487",
				"qq":           "123456",
				"avatar":       "",
				"created_by":   "0000000000",
				"gmt_create":   "2022-04-23 15:49:59",
				"gmt_modified": "2022-04-30 17:29:46",
			},
		},
	},
}

var UpdateMemberAvatarData = []APITestCase{}

var CreateMemberData = []APITestCase{
	{
		"success",
		Request{
			"POST",
			"/members/3000000000",
			"",
			gin.H{
				"alias":   "小稽",
				"name":    "滑小稽",
				"role":    "member_inactive",
				"section": "计算机233",
				"profile": "。。。",
				"phone":   "12352439487",
				"qq":      "123456",
			},
		},
		Response{
			200,
			gin.H{
				"member_id":    "3000000000",
				"alias":        "小稽",
				"name":         "滑小稽",
				"section":      "计算机233",
				"role":         "member_inactive",
				"profile":      "。。。",
				"phone":        "",
				"qq":           "123456",
				"avatar":       "",
				"created_by":   "2333333333",
				"gmt_create":   "2022-04-30 23:06:44",
				"gmt_modified": "2022-04-30 23:06:44",
			},
		},
	},
	{
		"invalid section",
		Request{
			"POST",
			"/members/3000000000",
			"",
			gin.H{
				"alias":   "小稽",
				"name":    "滑小稽",
				"section": "计算机23",
				"profile": "。。。",
				"phone":   "12352439487",
				"qq":      "123456",
			},
		},
		Response{
			422,
			nil,
		},
	},
}

var UpdateMemberBasicInfoData = []APITestCase{
	{
		"success",
		Request{
			"PATCH",
			"/members/2333333333",
			"",
			gin.H{
				"name":    "滑稽",
				"section": "计算机322",
				"role":    "admin",
			},
		},
		Response{
			200,
			gin.H{
				"member_id":    "2333333333",
				"alias":        "滑稽",
				"name":         "滑稽",
				"section":      "计算机322",
				"profile":      "relaxing",
				"phone":        "12356839487",
				"qq":           "123456",
				"avatar":       "",
				"created_by":   "0000000000",
				"gmt_create":   "2022-04-17T19:35:55.000Z",
				"gmt_modified": "2022-04-17T19:35:55.000Z",
				"role":         "admin",
			},
		},
	},
}

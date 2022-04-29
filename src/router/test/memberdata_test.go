package router_test

type H map[string]interface{}

type TestData struct {
	Method       string
	Url          string
	RequestBody  H
	Code         int
	ResponseBody H
}

var CreateMemberData = []TestData{
	{
		"POST",
		"/members/1231231231",
		H{
			"alias":   "123",
			"name":    "滑小稽",
			"section": "计算机123",
			"phone":   17523458765,
		},
		200,
		nil,
	},
}

var GetMemberData = []TestData{
	{
		"GET",
		"/members/2333333333",
		nil,
		200,
		H{
			"member_id":    "2333333333",
			"alias":        "huaji",
			"password":     "123456",
			"name":         "滑稽",
			"section":      "计算机233",
			"role":         "member",
			"profile":      "",
			"phone":        "",
			"qq":           "",
			"avatar":       "",
			"created_by":   "",
			"gmt_create":   "2022-04-23 15:49:59",
			"gmt_modified": "2022-04-23 15:50:01",
		},
	},
}

package model

type Member struct {
	MemberId    string `json:"member_id" `
	Alias       string `json:"alias"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Section     string `json:"section"`
	Profile     string `json:"profile"`
	Phone       string `json:"phone"`
	Qq          string `json:"qq"`
	Avatar      string `json:"avatar"`
	CreatedBy   string `json:"created_by"`
	GmtCreate   string `json:"gmt_create"`
	GmtModified string `json:"gmt_modified"`
	Role        string `json:"role"`
}

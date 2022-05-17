package model

type Client struct {
	ClientId int64 `json:"client_id"`
	Openid string `json:"openid"`
	GmtCreate string `json:"gmt_create"`
	GmtModified string `json:"gmt_modified"`
}
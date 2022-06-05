package model

type Client struct {
	ClientId    int64  `json:"client_id" db:"client_id"`
	OpenId      string `json:"openid"`
	GmtCreate   string `json:"gmt_create" db:"gmt_create"`
	GmtModified string `json:"gmt_modified" db:"gmt_modified"`
}

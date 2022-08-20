package model

type Client struct {
	ClientId    int64  `json:"clientId" db:"client_id"`
	OpenId      string `json:"openid"`
	GmtCreate   string `json:"gmtCreate" db:"gmt_create"`
	GmtModified string `json:"gmtModified" db:"gmt_modified"`
}

package model

type Member struct {
	MemberId    string `json:"member_id" db:"member_id"`
	Alias       string `json:"alias"`
	Password    string `json:"-"`
	Name        string `json:"name" `
	Section     string `json:"section" `
	Role        string `json:"role"`
	Profile     string `json:"profile"`
	Phone       string `json:"phone" `
	Qq          string `json:"qq" `
	Avatar      string `json:"avatar"`
	CreatedBy   string `json:"created_by" db:"created_by"`
	GmtCreate   string `json:"gmt_create" db:"gmt_create"`
	GmtModified string `json:"gmt_modified" db:"gmt_modified"`
}

type MemberRoleRelation struct {
	MemberId string `json:"member_id"`
	RoleId   int64  `json:"role_id"`
}
type Role struct {
	RoleId int64  `json:"role_id"`
	Role   string `json:"role"`
}

type PublicMember struct {
	MemberId    string `json:"member_id" db:"member_id"`
	Alias       string `json:"alias"`
	Role        string `json:"role"`
	Profile     string `json:"profile"`
	Avatar      string `json:"avatar"`
	CreatedBy   string `json:"created_by" db:"created_by"`
	GmtCreate   string `json:"gmt_create" db:"gmt_create"`
	GmtModified string `json:"gmt_modified" db:"gmt_modified"`
}

func CreatePublicMember(m Member) PublicMember {
	return PublicMember{
		MemberId:    m.MemberId,
		Alias:       m.Alias,
		Role:        m.Role,
		Profile:     m.Profile,
		Avatar:      m.Avatar,
		CreatedBy:   m.CreatedBy,
		GmtCreate:   m.GmtCreate,
		GmtModified: m.GmtModified,
	}
}

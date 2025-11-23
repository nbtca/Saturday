package model

import "database/sql"

type Member struct {
	MemberId    string `json:"memberId" db:"member_id"`
	LogtoId     string `json:"logtoId" db:"logto_id"`
	Alias       string `json:"alias"`
	Password    string `json:"-"`
	Name        string `json:"name" `
	Section     string `json:"section" `
	Role        string `json:"role"`
	Profile     string `json:"profile"`
	Phone       string `json:"phone" `
	QQ          string `json:"qq" `
	Avatar      string `json:"avatar"`
	CreatedBy   string `json:"createdBy" db:"created_by"`
	GmtCreate   string `json:"gmtCreate" db:"gmt_create"`
	GmtModified string `json:"gmtModified" db:"gmt_modified"`
}

type NullMember struct {
	MemberId    sql.NullString `json:"memberId" db:"member_id"`
	Alias       sql.NullString `json:"alias"`
	Password    sql.NullString `json:"-"`
	Name        sql.NullString `json:"name" `
	Section     sql.NullString `json:"section" `
	Role        sql.NullString `json:"role"`
	Profile     sql.NullString `json:"profile"`
	Phone       sql.NullString `json:"phone" `
	QQ          sql.NullString `json:"qq" `
	Avatar      sql.NullString `json:"avatar"`
	CreatedBy   sql.NullString `json:"createdBy" db:"created_by"`
	GmtCreate   sql.NullString `json:"gmtCreate" db:"gmt_create"`
	GmtModified sql.NullString `json:"gmtModified" db:"gmt_modified"`
}

func (nm NullMember) Member() *Member {
	if !nm.MemberId.Valid {
		return nil
	}
	return &Member{
		MemberId:    nm.MemberId.String,
		Alias:       nm.Alias.String,
		Role:        nm.Role.String,
		Profile:     nm.Profile.String,
		Section:     nm.Section.String,
		Phone:       nm.Phone.String,
		QQ:          nm.QQ.String,
		Avatar:      nm.Avatar.String,
		CreatedBy:   nm.CreatedBy.String,
		GmtCreate:   nm.GmtCreate.String,
		GmtModified: nm.GmtModified.String,
	}
}

func (nm NullMember) PublicMember() *PublicMember {
	if !nm.MemberId.Valid {
		return nil
	}
	return &PublicMember{
		MemberId:    nm.MemberId.String,
		Alias:       nm.Alias.String,
		Role:        nm.Role.String,
		Profile:     nm.Profile.String,
		Avatar:      nm.Avatar.String,
		CreatedBy:   nm.CreatedBy.String,
		GmtCreate:   nm.GmtCreate.String,
		GmtModified: nm.GmtModified.String,
	}
}

type MemberRoleRelation struct {
	MemberId string `json:"memberId"`
	RoleId   int64  `json:"role_id"`
}
type Role struct {
	RoleId int64  `json:"role_id"`
	Role   string `json:"role"`
}

type PublicMember struct {
	MemberId    string `json:"memberId" db:"member_id"`
	Alias       string `json:"alias"`
	Role        string `json:"role"`
	Profile     string `json:"profile"`
	Avatar      string `json:"avatar"`
	CreatedBy   string `json:"createdBy" db:"created_by"`
	GmtCreate   string `json:"gmtCreate" db:"gmt_create"`
	GmtModified string `json:"gmtModified" db:"gmt_modified"`
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

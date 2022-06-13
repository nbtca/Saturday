package repo

import (
	"database/sql"
	"saturday/model"
	"saturday/util"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

var memberFields = []string{"member_id", "alias", "password", "name", "section", "role", "profile",
	"phone", "qq", "avatar", "created_by", "gmt_create", "gmt_modified"}

func getMemberStatement() squirrel.SelectBuilder {
	return squirrel.Select("*").From("member_view")
}

func ExistMember(id string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) as count FROM member where member_id = ?", id)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func GetMemberById(id string) (model.Member, error) {
	statement, args, _ := getMemberStatement().Where(squirrel.Eq{"member_id": id}).ToSql()
	member := model.Member{}
	if err := db.Get(&member, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Member{}, nil
		}
		return model.Member{}, err
	}
	return member, nil
}

func GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	sql, args, _ := getMemberStatement().Offset(offset).Limit(limit).ToSql()
	members := []model.Member{}
	if err := db.Select(&members, sql, args...); err != nil {
		return []model.Member{}, err
	}
	return members, nil
}

func CreateMember(member *model.Member) error {
	member.GmtCreate = util.GetDate()
	member.GmtCreate = util.GetDate()
	sqlMember, argsMember, _ := squirrel.Insert("member").Columns(
		"member_id", "alias", "name", "section", "profile", "avatar",
		"phone", "qq", "created_by", "gmt_create", "gmt_modified").Values(
		member.MemberId, member.Alias, member.Name, member.Section,
		member.Profile, member.Avatar, member.Phone, member.QQ, member.CreatedBy,
		member.GmtCreate, member.GmtModified).ToSql()
	conn, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			conn.Rollback()
		}
	}()
	conn.Exec(sqlMember, argsMember...)
	SetMemberRole(member.MemberId, member.Role, conn)
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateMember(member model.Member) error {
	sql, args, _ := squirrel.Update("member").
		Set("alias", member.Alias).
		Set("name", member.Name).
		Set("section", member.Section).
		Set("profile", member.Profile).
		Set("phone", member.Phone).
		Set("qq", member.QQ).
		Set("gmt_modified", util.GetDate()).
		Where(squirrel.Eq{"member_id": member.MemberId}).ToSql()
	conn, err := db.Begin()
	defer func() {
		if err != nil {
			conn.Rollback()
		}
	}()
	if err != nil {
		return err
	}
	conn.Exec(sql, args...)
	err = SetMemberRole(member.MemberId, member.Role, conn)
	if err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

package repo

import (
	"database/sql"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"

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

func GetMemberIdByLogtoId(logtoId string) (sql.NullString, error) {
	var memberId sql.NullString
	err := db.Get(&memberId, "SELECT member_id FROM member where logto_id = ?", logtoId)
	if err != nil {
		return sql.NullString{}, err
	}
	return memberId, nil
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
	member.GmtModified = util.GetDate()
	sqlMember, argsMember, _ := squirrel.Insert("member").Columns(
		"member_id", "alias", "name", "section", "profile", "avatar",
		"phone", "qq", "created_by", "gmt_create", "gmt_modified").Values(
		member.MemberId, member.Alias, member.Name, member.Section,
		member.Profile, member.Avatar, member.Phone, member.QQ, member.CreatedBy,
		member.GmtCreate, member.GmtModified).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	defer util.RollbackOnErr(err, conn)
	if _, err = conn.Exec(sqlMember, argsMember...); err != nil {
		return err
	}
	if err = SetMemberRole(member.MemberId, member.Role, conn); err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateMember(member model.Member) error {
	sql, args, _ := squirrel.Update("member").
		Set("logto_id", member.LogtoId).
		Set("alias", member.Alias).
		Set("name", member.Name).
		Set("section", member.Section).
		Set("password", member.Password).
		Set("profile", member.Profile).
		Set("phone", member.Phone).
		Set("qq", member.QQ).
		Set("avatar", member.Avatar).
		Set("gmt_modified", util.GetDate()).
		Where(squirrel.Eq{"member_id": member.MemberId}).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	defer util.RollbackOnErr(err, conn)
	if _, err = conn.Exec(sql, args...); err != nil {
		return err
	}
	if err = SetMemberRole(member.MemberId, member.Role, conn); err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

package repo

import (
	"saturday/src/model"
	"saturday/src/util"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

var fields = []string{"member_id", "alias", "password", "name", "section", "role", "profile",
	"phone", "qq", "avatar", "created_by", "gmt_create", "gmt_modified"}

func getMemberStatement() squirrel.SelectBuilder {
	members := squirrel.Select(fields...).From("member")
	return members.LeftJoin("member_role_relation USING (member_id)").LeftJoin("role USING (role_id)")
}

func pagination(offset uint64, limit uint64) squirrel.SelectBuilder {
	return getMemberStatement().Offset(offset).Limit(limit)
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
	sql, args, _ := getMemberStatement().Where(squirrel.Eq{"member_id": id}).ToSql()
	member := model.Member{}
	if err := db.Get(&member, sql, args...); err != nil {
		return model.Member{}, err
	}
	return member, nil
}
func GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	sql, args, _ := pagination(offset, limit).ToSql()
	members := []model.Member{}
	if err := db.Select(&members, sql, args...); err != nil {
		return []model.Member{}, err
	}
	return members, nil
}

// TODO use 'row affected'
func CreateMember(member *model.Member) error {
	sqlMember, argsMember, _ := squirrel.Insert("member").Columns(
		"member_id", "alias", "name", "section", "profile",
		"phone", "qq", "created_by", "gmt_create", "gmt_modified").Values(
		member.MemberId, member.Alias, member.Name, member.Section,
		member.Profile, member.Phone, member.Qq, member.CreatedBy,
		util.GetDate(), util.GetDate()).ToSql()
	sqlRole, argsRole, _ := squirrel.Insert("member_role_relation").Columns("member_id", "role_id").Values(member.MemberId, 0).ToSql()
	conn, err := db.Begin()
	if err != nil {
		return err
	}
	conn.Exec(sqlMember, argsMember...)
	conn.Exec(sqlRole, argsRole...)
	if err = conn.Commit(); err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

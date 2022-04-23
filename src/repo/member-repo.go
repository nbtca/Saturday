package repo

import (
	"saturday/src/model"
	"saturday/src/util"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

// type Member struct{}

func getMemberStatement() squirrel.SelectBuilder {
	members := squirrel.Select(util.FieldsConstructor(model.Member{})).From("member")
	return members.LeftJoin("member_role_relation USING (member_id)").LeftJoin("role USING (role_id)")
}

func getMemberByIdStatement(id string) squirrel.SelectBuilder {
	return getMemberStatement().Where(squirrel.Eq{"member_id": id})
}

func pagination(offset uint64, limit uint64) squirrel.SelectBuilder {
	return getMemberStatement().Offset(offset).Limit(limit)
}

func GetMemberById(id string) (model.Member, error) {
	sql, args, _ := getMemberByIdStatement(id).ToSql()
	member := model.Member{}
	err := db.Get(&member, sql, args...)
	if err != nil {
		return model.Member{}, err
	}
	return member, nil
}

func GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	sql, args, _ := pagination(offset, limit).ToSql()
	members := []model.Member{}
	err := db.Select(&members, sql, args...)
	if err != nil {
		return []model.Member{}, err
	}
	return members, nil
}

func CreateMember(member *model.Member) error {
	sqlMember, argsMember, _ := squirrel.Insert("member").Columns("member_id", "alias", "name", "section", "profile", "phone", "qq", "created_by", "gmt_create", "gmt_modified").Values(member.MemberId, member.Alias, member.Name, member.Section, member.Profile, member.Phone, member.Qq, member.CreatedBy, time.Now().Format("2006-01-02 15:04:11"), time.Now().Format("2006-01-02 15:04:11")).ToSql()
	sqlRole, argsRole, _ := squirrel.Insert("member_role_relation").Columns("member_id", "role_id").Values(member.MemberId, 1).ToSql()
	conn, err := db.Begin()
	if err != nil {
		return err
	}
	conn.Exec(sqlMember, argsMember...)
	conn.Exec(sqlRole, argsRole...)
	err = conn.Commit()
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

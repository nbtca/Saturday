package repo

import (
	"gin-example/src/model"
	"gin-example/src/util"
	"log"
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

func GetMemberById(id string) model.Member {
	sql, args, _ := getMemberByIdStatement(id).ToSql()
	member := []model.Member{}
	err := db.Select(&member, sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	if len(member) != 0 {
		return member[0]
	} else {
		return model.Member{}
	}
}

func GetMembers(offset uint64, limit uint64) []model.Member {
	sql, args, _ := pagination(offset, limit).ToSql()
	log.Println(sql)
	members := []model.Member{}
	err := db.Select(&members, sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	return members
}

func CreateMember(member *model.Member) error {
	sqlMember, argsMember, _ := squirrel.Insert("member").Columns("member_id", "alias", "name", "section", "profile", "phone", "qq", "created_by", "gmt_create", "gmt_modified").Values(member.MemberId, member.Alias, member.Name, member.Section, member.Profile, member.Phone, member.Qq, member.CreatedBy, time.Now().Format("2006-01-02 15:04:11"), time.Now().Format("2006-01-02 15:04:11")).ToSql()
	sqlRole, argsRole, _ := squirrel.Insert("member_role_relation").Columns("member_id", "role_id").Values(member.MemberId, 1).ToSql()
	conn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	conn.Exec(sqlMember, argsMember...)
	conn.Exec(sqlRole, argsRole...)
	commitError := conn.Commit()
	if commitError != nil {
		log.Fatal(commitError)
		conn.Rollback()
		return commitError
	}
	return nil
}

package repo

import (
	"gin-example/src/model"
	"gin-example/src/util"
	"log"

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

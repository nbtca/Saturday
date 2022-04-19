package repo

import (
	model "gin-example/models"
	"gin-example/util"
	"log"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MemberRepo struct {
	DB *sqlx.DB
}

func (MemberRepo) getMemberStatement() squirrel.SelectBuilder {
	members := squirrel.Select(util.FieldsConstructor(model.Member{})).From("member")
	return members.LeftJoin("member_role_relation USING (member_id)").LeftJoin("role USING (role_id)")
}

func (repo *MemberRepo) getMemberByIdStatement(id string) squirrel.SelectBuilder {
	return repo.getMemberStatement().Where(squirrel.Eq{"member_id": id})
}

func (repo *MemberRepo) pagination(offset uint64, limit uint64) squirrel.SelectBuilder {
	return repo.getMemberStatement().Offset(offset).Limit(limit)
}

func (repo *MemberRepo) GetMemberById(id string) model.Member {
	sql, args, _ := repo.getMemberByIdStatement(id).ToSql()
	member := []model.Member{}
	err := repo.DB.Select(&member, sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	if len(member) != 0 {
		return member[0]
	} else {
		return model.Member{}
	}
}

func (repo *MemberRepo) GetMembers(offset uint64, limit uint64) []model.Member {
	sql, args, _ := repo.pagination(offset, limit).ToSql()
	log.Println(sql)
	members := []model.Member{}
	err := repo.DB.Select(&members, sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	return members
}

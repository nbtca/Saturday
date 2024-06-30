package repo

import (
	"database/sql"
	"time"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"

	"github.com/Masterminds/squirrel"
)

func GetClientByOpenId(openId string) (model.Client, error) {
	statement, args, _ := sq.Select("*").From("client").Where(squirrel.Eq{"openid": openId}).ToSql()
	client := model.Client{}
	if err := db.Get(&client, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Client{}, nil
		}
		return model.Client{}, nil
	}
	return client, nil
}

func CreateClient(client *model.Client) error {
	client.GmtCreate = util.GetDate()
	client.GmtModified = util.GetDate()
	sql, args, _ := sq.Insert("client").Columns("openid", "gmt_create", "gmt_modified").
		Values(client.OpenId, time.Now(), time.Now()).ToSql()
	res, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}
	client.ClientId, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

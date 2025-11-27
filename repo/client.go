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

func GetClientByLogtoId(logtoId string) (model.Client, error) {
	statement, args, _ := sq.Select("*").From("client").Where(squirrel.Eq{"logto_id": logtoId}).ToSql()
	client := model.Client{}
	if err := db.Get(&client, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Client{}, nil
		}
		return model.Client{}, nil
	}
	return client, nil
}

func GetClientById(clientId int64) (model.Client, error) {
	statement, args, _ := sq.Select("*").From("client").Where(squirrel.Eq{"client_id": clientId}).ToSql()
	client := model.Client{}
	if err := db.Get(&client, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Client{}, nil
		}
		return model.Client{}, err
	}
	return client, nil
}

func CreateClient(client *model.Client) error {
	client.GmtCreate = util.GetDate()
	client.GmtModified = util.GetDate()
	sql, args, _ := sq.Insert("client").Columns("openid", "logto_id", "gmt_create", "gmt_modified").
		Values(client.OpenId, client.LogtoId, time.Now(), time.Now()).ToSql()
	var id int64
	err := db.QueryRow(sql+"RETURNING client_id", args...).Scan(&id)
	if err != nil {
		return err
	}
	client.ClientId = id
	return nil
}

package repo

import (
	"os"
	"saturday/util"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() {
	var err error
	db, err = sqlx.Connect("mysql", os.Getenv("DB_URL"))
	if err != nil {
		util.Logger.Fatal(err)
	}
}

func SetDB(dbx *sqlx.DB) {
	db = dbx
}

func CloseDB() {
	db.Close()
}

package repo

import (
	"os"
	"time"

	"github.com/nbtca/saturday/util"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() {
	var err error
	db, err = sqlx.Connect("mysql", os.Getenv("DB_URL"))

	if err != nil {
		util.Logger.Fatal(err)
	}
	db.SetMaxOpenConns(1000)               // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)                 // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Minute * 5) // 0, connections are reused forever.
}

func SetDB(dbx *sqlx.DB) {
	db = dbx
}

func CloseDB() {
	db.Close()
}

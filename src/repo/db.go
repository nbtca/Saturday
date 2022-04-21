package repo

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitDB() {
	var err error
	db, err = sqlx.Connect("mysql", "root:password@/saturday_dev")
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

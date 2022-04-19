package util

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func (DBWrapper *DBWrapper) GetConnection() {
	db, err := sqlx.Connect("mysql", "root:password@/saturday_dev")
	if err != nil {
		log.Fatal(err)
	}
	DBWrapper.DB = db
}

func (DBWrapper *DBWrapper) CloseConnection() {
	DBWrapper.DB.Close()
}

func (DBWrapper *DBWrapper) GetDB() *sqlx.DB {
	return DBWrapper.DB
}

type DBWrapper struct {
	DB *sqlx.DB
}

var DB = &DBWrapper{}

package repo

import (
	"gin-example/src/util"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		util.Logger().Fatal("Error loading .env file")
	}

	db, err = sqlx.Connect("mysql", os.Getenv("DB_URL"))
	if err != nil {
		util.Logger().Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

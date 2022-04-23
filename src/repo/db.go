package repo

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

func InitDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err = sqlx.Connect("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

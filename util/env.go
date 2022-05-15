package util

import (
	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		Logger.Fatal("Error loading .env file")
	}
}

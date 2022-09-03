package main

import (
	"saturday/repo"
	"saturday/router"
	"saturday/util"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		util.Logger.Fatal("Error loading .env file")
	}

	util.InitValidator()
	util.InitDialer()

	repo.InitDB()
	defer repo.CloseDB()

	r := router.SetupRouter()
	r.Run(":8080")

	util.Logger.Info("Starting server...")
}

package main

import (
	"log"
	"os"

	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/util"

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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r.Run(":" + port)

	util.Logger.Info("Starting server at %v...", port)
}

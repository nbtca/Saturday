package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/util"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		util.Logger.Warning("Error loading .env file")
		log.Println(err)
	}

	util.InitValidator()
	util.InitDialer()

	repo.InitDB()
	defer repo.CloseDB()

	gin.DefaultWriter = util.Logger.Out
	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r.Run(":" + port)

	util.Logger.Infof("Starting server at %s...", port)
}

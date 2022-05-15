package main

import (
	"saturday/repo"
	"saturday/router"
	"saturday/util"
)

func main() {
	util.InitEnv()
	util.InitValidator()

	repo.InitDB()
	defer repo.CloseDB()

	r := router.SetupRouter()
	r.Run(":8080")

	util.Logger.Info("Starting server...")
}

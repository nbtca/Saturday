package main

import (
	"saturday/src/repo"
	"saturday/src/router"
	"saturday/src/util"
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

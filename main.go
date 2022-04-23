package main

import (
	"gin-example/src/repo"
	"gin-example/src/router"
	"gin-example/src/util"
)

func main() {
	repo.InitDB()
	defer repo.CloseDB()
	r := router.SetupRouter()
	r.Run(":8080")
	util.Logger().Info("Starting server...")
}

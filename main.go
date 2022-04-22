package main

import (
	"gin-example/src/repo"
	"gin-example/src/router"
	"gin-example/util"
)

func main() {
	util.Logger().Info("Starting server...")
	repo.InitDB()
	defer repo.CloseDB()
	router.InitRouter()
}

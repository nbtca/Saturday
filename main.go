package main

import (
	_ "github.com/go-sql-driver/mysql"

	"gin-example/src/repo"
	"gin-example/src/router"
)

func main() {
	repo.InitDB()
	defer repo.CloseDB()
	router.InitRouter()
}

package main

import (
	"database/sql"

	structMaker "github.com/Tsmwhite/structMaker/bin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConfigString := "root:password@tcp(127.0.0.1:3306)/saturday_dev?charset=utf8"
	db, err := sql.Open("mysql", dbConfigString)
	if err == nil {
		// @1
		//err = structMaker.Run(db, "example", structMaker.EgMySql)

		// @2
		loader := structMaker.NewMysql(db, "saturday_dev")
		err = structMaker.New().SetSourceDB(loader).MakeFile()

		// @3
		// err = structMaker.New().SetSourceDB(loader).SetOutput("models2", true).MakeFile()
		// fmt.Println(err)
	}
}

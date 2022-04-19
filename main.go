package main

import (
	_ "github.com/go-sql-driver/mysql"

	"gin-example/src/controller"
	"gin-example/src/repo"
	"gin-example/src/service"
	"gin-example/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	util.DB.GetConnection()
	defer util.DB.CloseConnection()
	repo := &repo.MemberRepo{
		DB: util.DB.GetDB(),
	}
	memberService := &service.MemberService{
		Repo: repo,
	}
	memberController := &controller.MemberController{
		MemberService: memberService,
	}
	r.GET("/:MemberId", memberController.GetMemberById)
	r.GET("/", memberController.GetByPage)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

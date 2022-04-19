package controller

import (
	"gin-example/src/service"
	"gin-example/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MemberController struct {
	MemberService *service.MemberService
}

func (controller *MemberController) GetMemberById(c *gin.Context) {
	validate := validator.New()
	memberId := c.Param("MemberId")
	errs := validate.Var(memberId, "required,len=10,numeric")
	if errs != nil {
		log.Println(errs)
	}
	member := controller.MemberService.GetMemberById(c.Param("MemberId"))
	c.JSON(200, member)
}

func (controller *MemberController) GetByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c)
	if err != nil {
		log.Println(err)
	}
	members := controller.MemberService.GetMembers(offset, limit)
	c.JSON(200, members)
}

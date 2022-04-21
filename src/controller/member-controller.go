package controller

import (
	"gin-example/src/service"
	"gin-example/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MemberController struct {
}

func (controller *MemberController) GetMemberById(c *gin.Context) {
	validate := validator.New()
	memberId := c.Param("MemberId")
	errs := validate.Var(memberId, "required,len=10,numeric")
	if errs != nil {
		log.Println(errs)
	}
	member, err := service.MemberServiceApp.GetMemberById(c.Param("MemberId"))
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "path fail"})
		c.Error(err)
		return
	}
	c.JSON(200, member)

}

func (controller *MemberController) GetByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c)
	if err != nil {
		log.Println(err)
	}
	members := service.MemberServiceApp.GetMembers(offset, limit)
	c.JSON(200, members)
}

var MemberControllerApp = new(MemberController)

package router

import (
	"gin-example/src/model"
	"gin-example/src/model/dto"
	"gin-example/src/service"
	"gin-example/src/util"
	"log"

	"github.com/gin-gonic/gin"
)

type MemberRouter struct {
}

func (MemberRouter) GetMemberById(c *gin.Context) {
	member, err := service.MemberServiceApp.GetMemberById(c.Param("MemberId"))
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) GetByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c) // TODO use validator
	if err != nil {
		c.Error(err)
	}
	members := service.MemberServiceApp.GetMembers(offset, limit)
	c.JSON(200, members)
}

func (MemberRouter) CreateToken(c *gin.Context) {
	CreateMemberTokenReq := &dto.CreateMemberTokenReq{}
	err := util.GetBody(c, CreateMemberTokenReq)
	log.Println(CreateMemberTokenReq)
	if util.CheckError(c, err) {
		return
	}
	res, err := service.MemberServiceApp.CreateToken(service.MemberAccount{
		MemberId: CreateMemberTokenReq.MemberId,
		Password: CreateMemberTokenReq.Password,
	})
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, res)
}

func (MemberRouter) Create(c *gin.Context) {
	CreateMemberReq := &model.Member{}
	err := util.GetBody(c, CreateMemberReq)
	CreateMemberReq.MemberId = c.Param("MemberId")
	if util.CheckError(c, err) {
		return
	}
	res, err := service.MemberServiceApp.CreateMember(CreateMemberReq)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, res)
}

var MemberRouterApp = new(MemberRouter)

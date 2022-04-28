package router

import (
	"saturday/src/model"
	"saturday/src/model/dto"
	"saturday/src/service"
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

type MemberRouter struct {
}

func (MemberRouter) GetMemberById(c *gin.Context) {
	memberId := &dto.MemberId{}
	err := util.BindAll(c, memberId)
	if util.CheckError(c, err) {
		return
	}
	member, err := service.MemberServiceApp.GetMemberById(memberId.MemberId)
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
	members, err := service.MemberServiceApp.GetMembers(offset, limit)
	if err != nil {
		c.Error(err)
	}
	c.JSON(200, members)
}

func (MemberRouter) CreateToken(c *gin.Context) {
	CreateMemberTokenReq := &dto.CreateMemberTokenReq{}
	err := util.BindAll(c, CreateMemberTokenReq)
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
	CreateMemberReq := &dto.CreateMemberReq{}
	err := util.BindAll(c, CreateMemberReq)
	if util.CheckError(c, err) {
		return
	}
	member := &model.Member{}
	err = util.SwapObject(CreateMemberReq, member)
	if err != nil {
		util.Logger.Error(err)
		c.JSON(500, err)
		return
	}
	res, err := service.MemberServiceApp.CreateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, res)
}

var MemberRouterApp = new(MemberRouter)

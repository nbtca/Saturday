package router

import (
	"net/http"
	"saturday/src/model"
	"saturday/src/model/dto"
	"saturday/src/service"
	"saturday/src/util"

	"github.com/gin-gonic/gin"
)

type MemberRouter struct {
}

func (MemberRouter) GetPublicMemberById(c *gin.Context) {
	memberId := &dto.MemberId{}
	if err := util.BindAll(c, memberId); util.CheckError(c, err) {
		return
	}
	member, err := service.MemberServiceApp.GetPublicMemberById(memberId.MemberId)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) GetPublicMemberByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c) // TODO use validator
	if err != nil {
		c.Error(err)
		return
	}
	members, err := service.MemberServiceApp.GetPublicMembers(offset, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, members)
}

func (MemberRouter) GetMemberById(c *gin.Context) {
	memberId := c.GetString("id")
	member, err := service.MemberServiceApp.GetMemberById(memberId)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) CreateToken(c *gin.Context) {
	req := &dto.CreateMemberTokenReq{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	member, err := service.MemberServiceApp.GetMemberById(req.MemberId)
	if util.CheckError(c, err) {
		return
	}
	if member.Password != req.Password {
		c.AbortWithStatusJSON(util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			Build())
	}
	token, err := service.MemberServiceApp.CreateToken(member)
	if util.CheckError(c, err) {
		return
	}
	res := dto.CreateMemberTokenResponse{
		Member: member,
		Token:  token,
	}
	c.JSON(200, res)
}

func (MemberRouter) Create(c *gin.Context) {
	req := &dto.CreateMemberReq{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	member := &model.Member{
		MemberId:  req.MemberId,
		Alias:     req.Alias,
		Name:      req.Name,
		Section:   req.Section,
		Profile:   req.Profile,
		Qq:        req.Qq,
		CreatedBy: c.GetString("id"),
	}
	res, err := service.MemberServiceApp.CreateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, res)
}

func (MemberRouter) CreateMany(c *gin.Context) {
	//TODO not implemented
}

func (MemberRouter) Activate(c *gin.Context) {
	//TODO not implemented
}

func (MemberRouter) Update(c *gin.Context) {
	//TODO not implemented
}

func (MemberRouter) UpdateBasic(c *gin.Context) {
	//TODO not implemented
}

func (MemberRouter) UpdateAvatar(c *gin.Context) {
	//TODO not implemented
}

var MemberRouterApp = new(MemberRouter)

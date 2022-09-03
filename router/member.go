package router

import (
	"net/http"
	"saturday/model"
	"saturday/model/dto"
	"saturday/service"
	"saturday/util"

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
	offset, limit, err := util.GetPaginationQuery(c)
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
	req := &dto.CreateMemberTokenRequest{}
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
		return
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
	req := &dto.CreateMemberRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	member := &model.Member{
		MemberId:  req.MemberId,
		Alias:     req.Alias,
		Name:      req.Name,
		Section:   req.Section,
		Avatar:    req.Avatar,
		Profile:   req.Profile,
		QQ:        req.QQ,
		Phone:     req.Phone,
		Role:      req.Role,
		CreatedBy: c.GetString("id"),
	}
	err := service.MemberServiceApp.CreateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) CreateMany(c *gin.Context) {
	//TODO not implemented
}

func (MemberRouter) Activate(c *gin.Context) {
	req := &dto.ActivateMemberRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	req.MemberId = c.GetString("id")
	member, err := service.MemberServiceApp.GetMemberById(req.MemberId)
	if util.CheckError(c, err) {
		return
	}
	member.Password = req.Password
	if req.Alias != "" {
		member.Alias = req.Alias
	}
	if req.Phone != "" {
		member.Phone = req.Phone
	}
	if req.QQ != "" {
		member.QQ = req.QQ
	}
	if req.Profile != "" {
		member.Profile = req.Profile
	}
	err = service.MemberServiceApp.ActivateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) Update(c *gin.Context) {
	req := &dto.UpdateMemberRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	req.MemberId = c.GetString("id")
	member, err := service.MemberServiceApp.GetMemberById(req.MemberId)
	if util.CheckError(c, err) {
		return
	}
	// TODO simplify
	if req.Alias != "" {
		member.Alias = req.Alias
	}
	if req.Phone != "" {
		member.Phone = req.Phone
	}
	if req.QQ != "" {
		member.QQ = req.QQ
	}
	if req.Avatar != "" {
		member.Avatar = req.Avatar
	}
	if req.Profile != "" {
		member.Profile = req.Profile
	}
	if req.Password != "" {
		member.Password = req.Password
	}
	err = service.MemberServiceApp.UpdateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) UpdateBasic(c *gin.Context) {
	req := &dto.UpdateMemberBasicRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	member, err := service.MemberServiceApp.GetMemberById(req.MemberId)
	if util.CheckError(c, err) {
		return
	}
	if req.Name != "" {
		member.Name = req.Name
	}
	if req.Section != "" {
		member.Section = req.Section
	}
	if req.Role != "" {
		member.Role = req.Role
	}
	err = service.MemberServiceApp.UpdateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) UpdateAvatar(c *gin.Context) {
	memberId := c.GetString("id")
	req := &dto.UpdateAvatarRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}
	member, err := service.MemberServiceApp.GetMemberById(memberId)
	if util.CheckError(c, err) {
		return
	}
	member.Avatar = req.Url
	err = service.MemberServiceApp.UpdateMember(member)
	if util.CheckError(c, err) {
		return
	}
	c.JSON(200, member)

}

var MemberRouterApp = new(MemberRouter)

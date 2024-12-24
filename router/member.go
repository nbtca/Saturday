package router

import (
	"log"
	"net/http"
	"os"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type MemberRouter struct{}

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

func (MemberRouter) GetMemberByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c)
	if err != nil {
		c.Error(err)
		return
	}
	members, err := service.MemberServiceApp.GetMembers(offset, limit)
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
			AddDetailError("member", "password", "invalid password").
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

func (MemberRouter) CreateTokenViaLogtoToken(c *gin.Context) {
	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))

	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if util.CheckError(c, err) {
		return
	}
	accessToken := res["access_token"].(string)

	auth := c.GetHeader("Authorization")
	user, err := service.LogtoServiceApp.FetchUserByToken(auth, accessToken)
	if util.CheckError(c, err) {
		return
	}
	invalidTokenError := util.
		MakeServiceError(http.StatusUnprocessableEntity).
		AddDetailError("member", "logto token", "invalid token")
	if user["id"] == nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: id missing").Build())
		return
	}
	logto_id, ok := user["id"].(string)
	if !ok {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: failed at getting id").Build())
		return
	}
	memberId, err := repo.GetMemberIdByLogtoId(logto_id)
	if err != nil || !memberId.Valid {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: member not found").Build())
		return
	}

	member, err := service.MemberServiceApp.GetMemberById(memberId.String)
	if util.CheckError(c, err) {
		return
	}
	logto_roles, err := service.LogtoServiceApp.FetchUserRole(logto_id, accessToken)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	member.Role = service.MemberServiceApp.MapLogtoUserRole(logto_roles)
	err = service.MemberServiceApp.UpdateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("error at syncing member role" + err.Error())
	}

	t, err := service.MemberServiceApp.CreateToken(member)
	if util.CheckError(c, err) {
		return
	}

	patchLogtoUserRequest := dto.PatchLogtoUserRequest{}

	logtoName, _ := user["name"].(string)
	if member.Alias != "" && logtoName == "" {
		patchLogtoUserRequest.Name = member.Alias
	}
	logtoAvatar, _ := user["avatar"].(string)
	if member.Avatar != "" && logtoAvatar == "" {
		patchLogtoUserRequest.Avatar = member.Avatar
	}

	_, err = service.LogtoServiceApp.PatchUserById(logto_id, patchLogtoUserRequest, accessToken)
	if err != nil {
		log.Println(err)
	}
	response := dto.CreateMemberTokenResponse{
		Member: member,
		Token:  t,
	}
	c.JSON(200, response)
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

func (MemberRouter) CreateWithLogto(c *gin.Context) {
	req := &dto.CreateMemberWithLogtoRequest{}
	if err := util.BindAll(c, req); util.CheckError(c, err) {
		return
	}

	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if util.CheckError(c, err) {
		return
	}
	accessToken := res["access_token"].(string)

	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))
	auth := c.GetHeader("Authorization")
	user, err := service.LogtoServiceApp.FetchUserByToken(auth, accessToken)
	if util.CheckError(c, err) {
		return
	}
	invalidTokenError := util.
		MakeServiceError(http.StatusUnprocessableEntity).
		AddDetailError("member", "logto token", "invalid token")
	if user["id"] == nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: id missing").Build())
		return
	}
	logtoId, ok := user["id"].(string)
	if !ok {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: failed at getting id").Build())
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
		Role:      "member",
		CreatedBy: req.MemberId,
		LogtoId:   logtoId,
	}
	err = service.MemberServiceApp.CreateMember(member)
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

func (MemberRouter) BindMemberLogtoId(c *gin.Context) {
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
			AddDetailError("member", "password", "invalid password").
			Build())
		return
	}

	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))

	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if util.CheckError(c, err) {
		return
	}
	accessToken := res["access_token"].(string)

	auth := c.GetHeader("Authorization")
	user, err := service.LogtoServiceApp.FetchUserByToken(auth, accessToken)
	if util.CheckError(c, err) {
		return
	}
	invalidTokenError := util.
		MakeServiceError(http.StatusUnprocessableEntity).
		AddDetailError("member", "logto token", "invalid token")
	if user["id"] == nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: id missing").Build())
		return
	}
	logtoId, ok := user["id"].(string)
	if !ok {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token: failed at getting id").Build())
		return
	}
	if member.LogtoId != "" {
		c.AbortWithStatusJSON(util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "logtoId", "already bound").
			Build())
		return
	}

	member.LogtoId = logtoId
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

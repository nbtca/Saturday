package router

import (
	"net/http"
	"net/url"
	"os"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
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

	auth := c.GetHeader("Authorization")
	// validate JWT and signature
	jwksURL, err := url.JoinPath(os.Getenv("LOGTO_ENDPOINT"), "/oidc/jwks")
	if util.CheckError(c, err) {
		return
	}
	_, claims, error := util.ParseTokenWithJWKS(jwksURL, auth)
	invalidTokenError := util.
		MakeServiceError(http.StatusUnprocessableEntity).
		AddDetailError("member", "logto token", "invalid token")
		// SetMessage("Invalid token"+error.Error()).
	if error != nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token" + error.Error()).Build())
		return
	}
	// check issuer
	// TODO move logto domain to config
	expectedIssuer, _ := url.JoinPath(os.Getenv("LOGTO_ENDPOINT"), "/oidc")
	if claims.Issuer != expectedIssuer {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token, invalid issuer").Build())
		return
	}
	// check audience
	// TODO move current resource indicator to config
	// expectedAudience := "https://api.nbtca.space/v2"
	// if claims.Audience != expectedAudience {
	// 	c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token").Build())
	// 	return
	// }
	// TODO check scope

	userId := claims.Subject
	res, err := service.LogtoServiceApp.FetchLogtoToken("https://default.logto.app/api", "all")
	if err != nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token").Build())
		return
	}
	accessToken := res["access_token"].(string)
	user, err := service.LogtoServiceApp.FetchUserById(userId, "Bearer "+accessToken)
	if err != nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token").Build())
		return
	}

	if user["customData"] == nil {
		c.AbortWithStatusJSON(invalidTokenError.SetMessage("Invalid token").Build())
	}

	customData, ok := user["customData"].(map[string]interface{})
	if !ok {
		c.AbortWithStatusJSON(util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "logto token", "invalid token").
			Build())
		return
	}
	memberId, ok := customData["memberId"].(string)
	if !ok {
		c.AbortWithStatusJSON(util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "logto token", "invalid token").
			Build())
		return
	}
	member, err := service.MemberServiceApp.GetMemberById(memberId)
	if util.CheckError(c, err) {
		return
	}
	t, err := service.MemberServiceApp.CreateToken(member)
	if util.CheckError(c, err) {
		return
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

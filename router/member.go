package router

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type MemberRouter struct{}

func (MemberRouter) GetPublicMemberById(ctx context.Context, input *struct {
	MemberId string `path:"MemberId" maxLength:"10" example:"2333333333" doc:"Name to greet"`
}) (*util.CommonResponse[model.PublicMember], error) {
	member, err := service.MemberServiceApp.GetPublicMemberById(input.MemberId)
	if err != nil {
		return nil, huma.NewError(http.StatusUnprocessableEntity, err.Error())
	}
	return util.MakeCommonResponse(member), nil
}

func (MemberRouter) GetPublicMemberByPage(ctx context.Context, input *struct {
	dto.PageRequest
}) (*util.CommonResponse[[]model.PublicMember], error) {
	members, err := service.MemberServiceApp.GetPublicMembers(input.Offset, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusUnprocessableEntity, err.Error())
	}
	return util.MakeCommonResponse(members), nil
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

func (MemberRouter) CreateToken(ctx context.Context, input *struct {
	MemberId string `path:"MemberId" maxLength:"10" example:"2333333333" doc:"Member Id"`
	Body     struct {
		Password string `json:"password"`
	}
}) (*util.CommonResponse[dto.CreateMemberTokenResponse], error) {
	member, err := service.MemberServiceApp.GetMemberById(input.MemberId)
	if err != nil {
		return nil, huma.NewError(http.StatusUnprocessableEntity, err.Error())
	}
	if member.Password != input.Body.Password {
		return nil, huma.NewError(http.StatusUnprocessableEntity, "Invalid password")
	}
	token, err := service.MemberServiceApp.CreateToken(member)
	if err != nil {
		return nil, huma.NewError(http.StatusUnprocessableEntity, err.Error())
	}
	return util.MakeCommonResponse(dto.CreateMemberTokenResponse{
		Member: member,
		Token:  token,
	}), nil
}

func (MemberRouter) CreateTokenViaLogtoToken(c context.Context, input *struct {
	MemberId      string `path:"MemberId" maxLength:"10" example:"2333333333" doc:"Member Id"`
	Authorization string `header:"Authorization"`
}) (*util.CommonResponse[dto.CreateMemberTokenResponse], error) {
	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))

	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	accessToken := res["access_token"].(string)

	user, err := service.LogtoServiceApp.FetchUserByToken(input.Authorization, accessToken)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if user["id"] == nil {
		return nil, huma.Error422UnprocessableEntity("Invalid token: id missing")
	}
	logto_id, ok := user["id"].(string)
	if !ok {
		return nil, huma.Error422UnprocessableEntity("Invalid token: failed at getting id")
	}
	memberId, err := repo.GetMemberIdByLogtoId(logto_id)
	if err != nil || !memberId.Valid {
		return nil, huma.Error422UnprocessableEntity("Invalid token: member not found")
	}

	member, err := service.MemberServiceApp.GetMemberById(memberId.String)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	t, err := service.MemberServiceApp.CreateToken(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
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
	return util.MakeCommonResponse(dto.CreateMemberTokenResponse{
		Member: member,
		Token:  t,
	}), nil
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

func (MemberRouter) BindMemberLogtoId(c context.Context, input *struct {
	MemberId      string `path:"MemberId" maxLength:"10" example:"2333333333" doc:"Member Id"`
	Authorization string `header:"Authorization"`
	Body          struct {
		Password string `json:"password" binding:""`
	}
}) (*util.CommonResponse[model.Member], error) {
	member, err := service.MemberServiceApp.GetMemberById(input.MemberId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if member.Password != input.Body.Password {
		return nil, huma.NewError(http.StatusUnprocessableEntity, "Invalid password")
	}

	service.LogtoServiceApp = service.MakeLogtoService(os.Getenv("LOGTO_ENDPOINT"))

	res, err := service.LogtoServiceApp.FetchLogtoToken(service.DefaultLogtoResource, "all")
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	accessToken := res["access_token"].(string)

	user, err := service.LogtoServiceApp.FetchUserByToken(input.Authorization, accessToken)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if user["id"] == nil {
		return nil, huma.Error422UnprocessableEntity("Invalid token: id missing")
	}
	logtoId, ok := user["id"].(string)
	if !ok {
		return nil, huma.Error422UnprocessableEntity("Invalid token: failed at getting id")
	}
	if member.LogtoId != "" {
		return nil, huma.Error422UnprocessableEntity("Validation Failed: member logtoId already bound")
	}

	member.LogtoId = logtoId
	err = service.MemberServiceApp.UpdateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(member), nil
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

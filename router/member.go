package router

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
	"github.com/spf13/viper"

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
	Authorization string `header:"Authorization"`
}) (*util.CommonResponse[dto.CreateMemberTokenResponse], error) {
	user, err := service.LogtoServiceApp.FetchUserByToken(input.Authorization)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if user.Id == "" {
		return nil, huma.Error422UnprocessableEntity("Invalid token: id missing")
	}
	memberId, err := repo.GetMemberIdByLogtoId(user.Id)
	if err != nil || !memberId.Valid {
		return nil, huma.Error422UnprocessableEntity("Invalid token: member not found")
	}

	member, err := service.MemberServiceApp.GetMemberById(memberId.String)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	logto_roles, err := service.LogtoServiceApp.FetchUserRole(user.Id)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	mappedRole := service.MemberServiceApp.MapLogtoUserRole(logto_roles)
	if mappedRole != member.Role && mappedRole != "" {
		member.Role = mappedRole
		err = service.MemberServiceApp.UpdateMember(member)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("error at syncing member role" + err.Error())
		}
	}

	t, err := service.MemberServiceApp.CreateToken(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	patchLogtoUserRequest := dto.PatchLogtoUserRequest{}

	if member.Alias != "" && user.Name == "" {
		patchLogtoUserRequest.Name = member.Alias
	}
	if member.Avatar != "" && user.Avatar == "" {
		patchLogtoUserRequest.Avatar = member.Avatar
	}

	_, err = service.LogtoServiceApp.PatchUserById(user.Id, patchLogtoUserRequest)
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

func (MemberRouter) CreateWithLogto(c context.Context, input *struct {
	MemberId string `path:"MemberId" maxLength:"10" example:"2333333333" doc:"Member Id"`
	LogtoId  string `json:"logtoId" doc:"Logto Id"`
	Name     string `json:"name" minLength:"2" maxLength:"4" doc:"Name"`
	Section  string `json:"section"`
	Alias    string `json:"alias" maxLength:"20"`
	Avatar   string `json:"avatar" maxLength:"255"`
	Profile  string `json:"profile" maxLength:"1000"`
	Phone    string `json:"phone"`
	QQ       string `json:"qq" minLength:"5" maxLength:"20"`
	Auth     string `header:"Authorization"`
}) (*util.CommonResponse[model.Member], error) {
	service.LogtoServiceApp = service.MakeLogtoService(viper.GetString("LOGTO_ENDPOINT"))
	user, err := service.LogtoServiceApp.FetchUserByToken(input.Auth)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if user.Id == "" {
		return nil, huma.Error422UnprocessableEntity("Invalid token: id missing")
	}

	member := &model.Member{
		MemberId:  input.MemberId,
		Alias:     input.Alias,
		Name:      input.Name,
		Section:   input.Section,
		Avatar:    input.Avatar,
		Profile:   input.Profile,
		QQ:        input.QQ,
		Phone:     input.Phone,
		Role:      "member",
		CreatedBy: input.MemberId,
		LogtoId:   user.Id,
	}
	if gh, ok := user.Identities["github"]; ok {
		member.GithubId = gh.UserId
	}

	if err = service.MemberServiceApp.CreateMember(member); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(*member), nil

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
		Password string `json:"password"`
	}
}) (*util.CommonResponse[model.Member], error) {
	member, err := service.MemberServiceApp.GetMemberById(input.MemberId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if member.Password != input.Body.Password {
		return nil, huma.NewError(http.StatusUnprocessableEntity, "Invalid password")
	}

	user, err := service.LogtoServiceApp.FetchUserByToken(input.Authorization)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if user.Id == "" {
		return nil, huma.Error422UnprocessableEntity("Invalid token: id missing")
	}
	if member.LogtoId != "" {
		return nil, huma.Error422UnprocessableEntity("Validation Failed: member logtoId already bound")
	}

	member.LogtoId = user.Id
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

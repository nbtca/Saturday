package router

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
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

func (MemberRouter) GetMemberByPage(ctx context.Context, input *GetMemberByPageInput) (*util.CommonResponse[[]model.Member], error) {
	_, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}
	members, err := service.MemberServiceApp.GetMembers(input.Offset, input.Limit)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(members), nil
}

func (MemberRouter) GetMemberById(ctx context.Context, input *GetMemberInput) (*util.CommonResponse[model.Member], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}
	member, err := service.MemberServiceApp.GetMemberById(auth.ID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(member), nil
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

func (MemberRouter) Create(ctx context.Context, input *CreateMemberInput) (*util.CommonResponse[model.Member], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}
	member := &model.Member{
		MemberId:  input.MemberId,
		Alias:     input.Body.Alias,
		Name:      input.Body.Name,
		Section:   input.Body.Section,
		Avatar:    input.Body.Avatar,
		Profile:   input.Body.Profile,
		QQ:        input.Body.QQ,
		Phone:     input.Body.Phone,
		Role:      input.Body.Role,
		CreatedBy: auth.ID,
	}
	err = service.MemberServiceApp.CreateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(*member), nil
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

func (MemberRouter) CreateMany(ctx context.Context, input *CreateManyMembersInput) (*util.CommonResponse[[]model.Member], error) {
	_, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}
	//TODO not implemented
	return nil, huma.Error501NotImplemented("CreateMany not implemented")
}

func (MemberRouter) Activate(ctx context.Context, input *ActivateMemberInput) (*util.CommonResponse[model.Member], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member_inactive", "admin_inactive")
	if err != nil {
		return nil, err
	}
	member, err := service.MemberServiceApp.GetMemberById(auth.ID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	member.Password = input.Body.Password
	if input.Body.Alias != "" {
		member.Alias = input.Body.Alias
	}
	if input.Body.Phone != "" {
		member.Phone = input.Body.Phone
	}
	if input.Body.QQ != "" {
		member.QQ = input.Body.QQ
	}
	if input.Body.Profile != "" {
		member.Profile = input.Body.Profile
	}
	err = service.MemberServiceApp.ActivateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(member), nil
}

func (MemberRouter) Update(ctx context.Context, input *UpdateMemberInput) (*util.CommonResponse[model.Member], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}
	member, err := service.MemberServiceApp.GetMemberById(auth.ID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	// TODO simplify
	if input.Body.Alias != "" {
		member.Alias = input.Body.Alias
	}
	if input.Body.Phone != "" {
		member.Phone = input.Body.Phone
	}
	if input.Body.QQ != "" {
		member.QQ = input.Body.QQ
	}
	if input.Body.Avatar != "" {
		member.Avatar = input.Body.Avatar
	}
	if input.Body.Profile != "" {
		member.Profile = input.Body.Profile
	}
	if input.Body.Password != "" {
		member.Password = input.Body.Password
	}
	err = service.MemberServiceApp.UpdateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(member), nil
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

func (MemberRouter) UpdateBasic(ctx context.Context, input *UpdateMemberBasicInput) (*util.CommonResponse[model.Member], error) {
	_, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}
	member, err := service.MemberServiceApp.GetMemberById(input.MemberId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if input.Body.Name != "" {
		member.Name = input.Body.Name
	}
	if input.Body.Section != "" {
		member.Section = input.Body.Section
	}
	if input.Body.Role != "" {
		member.Role = input.Body.Role
	}
	err = service.MemberServiceApp.UpdateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(member), nil
}

func (MemberRouter) UpdateAvatar(ctx context.Context, input *UpdateMemberAvatarInput) (*util.CommonResponse[model.Member], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}
	member, err := service.MemberServiceApp.GetMemberById(auth.ID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	member.Avatar = input.Body.Avatar
	err = service.MemberServiceApp.UpdateMember(member)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(member), nil
}

var MemberRouterApp = new(MemberRouter)

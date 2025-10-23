package service

import (
	"net/http"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"
)

type MemberService struct {
	memberRepo repo.MemberRepository
	roleRepo   repo.RoleRepository
	logtoService LogtoServiceInterface
}

// NewMemberService creates a new MemberService with injected dependencies
func NewMemberService(memberRepo repo.MemberRepository, roleRepo repo.RoleRepository, logtoService LogtoServiceInterface) MemberServiceInterface {
	return &MemberService{
		memberRepo: memberRepo,
		roleRepo: roleRepo,
		logtoService: logtoService,
	}
}

func (service *MemberService) GetMemberById(id string) (model.Member, error) {
	member, err := service.memberRepo.GetMemberById(id)
	if err != nil {
		return model.Member{}, err
	}
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed").
			AddDetailError("member", "memberId", "invalid memberId")
		return member, error
	} else {
		return member, nil
	}
}

func (service *MemberService) GetMemberByLogtoId(logtoId string) (model.Member, error) {
	member, err := service.memberRepo.GetMemberByLogtoId(logtoId)
	if err != nil {
		return model.Member{}, err
	}
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed").
			AddDetailError("member", "logtoId", "invalid logtoId")
		return member, error
	} else {
		return member, nil
	}
}

func (service *MemberService) GetMemberByGithubId(githubId string) (model.Member, error) {
	member, err := service.memberRepo.GetMemberByGithubId(githubId)
	if err != nil {
		return model.Member{}, err
	}
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed").
			AddDetailError("member", "logtoId", "invalid logtoId")
		return member, error
	} else {
		return member, nil
	}
}

func (service *MemberService) GetPublicMemberById(id string) (model.PublicMember, error) {
	member, err := service.GetMemberById(id)
	if err != nil {
		return model.PublicMember{}, err
	}
	return model.CreatePublicMember(member), nil
}

func (service *MemberService) GetPublicMembers(offset uint64, limit uint64) ([]model.PublicMember, error) {
	members, err := service.memberRepo.GetMembers(offset, limit)
	if err != nil {
		return nil, err
	}
	var publicMembers []model.PublicMember
	for _, member := range members {
		publicMembers = append(publicMembers, model.CreatePublicMember(member))
	}
	return publicMembers, nil
}

func (service *MemberService) GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	members, err := service.memberRepo.GetMembers(offset, limit)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (service *MemberService) CreateMember(member *model.Member) error {
	if member.Role != "admin" && member.Role != "member" {
		return util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "role", "invalid role")
	}
	exist, err := service.memberRepo.ExistMember(member.MemberId)
	if err != nil {
		return err
	}
	if exist {
		return util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
	}
	if member.Role == "member" {
		member.Role = "member_inactive"
	}
	if member.Role == "admin" {
		member.Role = "admin_inactive"
	}
	if err := service.memberRepo.CreateMember(member); err != nil {
		return err
	}
	return nil
}

func (service *MemberService) CreateToken(member model.Member) (string, error) {
	res, err := util.CreateToken(util.Payload{Who: member.MemberId, Member: member, Role: member.Role})
	return res, err
}

// func (service *MemberService) UpdateBasic(member model.Member) error {
// 	exist, err := service.roleRepo.ExistRole(member.Role)
// 	if err != nil {
// 		return err
// 	}
// 	if !exist {
// 		return util.
// 			MakeServiceError(http.StatusUnprocessableEntity).
// 			SetMessage("Validation Failed").
// 			AddDetailError("member", "role", "invalid role")
// 	}
// 	if err := service.memberRepo.UpdateMember(member); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (service MemberService) MapLogtoUserRole(roles []LogtoUserRole) string {
	role := ""
	for _, r := range roles {
		if r.Name == "Repair Admin" {
			role = "admin"
		} else if r.Name == "Repair Member" && role == "" {
			role = "member"
		}
	}
	return role
}

func (service *MemberService) UpdateMember(member model.Member) error {
	exist, err := service.roleRepo.ExistRole(member.Role)
	if err != nil {
		return err
	}
	if !exist {
		return util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "role", "invalid role")
	}
	if err := service.memberRepo.UpdateMember(member); err != nil {
		return err
	}
	return nil
}

func (service *MemberService) ActivateMember(member model.Member) error {
	if member.Role == "member_inactive" {
		member.Role = "member"
	}
	if member.Role == "admin_inactive" {
		member.Role = "admin"
	}
	if err := service.memberRepo.UpdateMember(member); err != nil {
		return err
	}
	return nil
}

// GetOrCreateMemberByLogtoId gets a member by Logto ID, handling the business logic
func (service *MemberService) GetOrCreateMemberByLogtoId(logtoId string) (model.Member, error) {
	memberId, err := service.memberRepo.GetMemberIdByLogtoId(logtoId)
	if err != nil || !memberId.Valid {
		return model.Member{}, util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "logtoId", "member not found for logto id")
	}
	
	return service.GetMemberById(memberId.String)
}

// AuthenticateAndSyncMember handles the complete authentication and sync flow from router
func (service *MemberService) AuthenticateAndSyncMember(logtoToken string) (model.Member, string, error) {
	// Validate Logto token and get user info
	user, err := service.logtoService.FetchUserByToken(logtoToken)
	if err != nil {
		return model.Member{}, "", err
	}
	if user.Id == "" {
		return model.Member{}, "", util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Invalid token: id missing")
	}

	// Get or create member
	member, err := service.GetOrCreateMemberByLogtoId(user.Id)
	if err != nil {
		return model.Member{}, "", err
	}

	// Sync role from Logto
	logtoRoles, err := service.logtoService.FetchUserRole(user.Id)
	if err != nil {
		return model.Member{}, "", err
	}

	mappedRole := service.MapLogtoUserRole(logtoRoles)
	if mappedRole != member.Role && mappedRole != "" {
		member.Role = mappedRole
		err = service.UpdateMember(member)
		if err != nil {
			return model.Member{}, "", util.MakeServiceError(http.StatusUnprocessableEntity).
				SetMessage("error at syncing member role: " + err.Error())
		}
	}

	// Sync profile to Logto if needed
	err = service.SyncMemberProfile(member, user)
	if err != nil {
		// Log error but don't fail the authentication
		// This is non-critical operation
	}

	// Create JWT token
	token, err := service.CreateToken(member)
	if err != nil {
		return model.Member{}, "", err
	}

	return member, token, nil
}

// SyncMemberProfile syncs member profile data to Logto user
func (service *MemberService) SyncMemberProfile(member model.Member, logtoUser *FetchLogtoUsersResponse) error {
	patchRequest := dto.PatchLogtoUserRequest{}
	needsUpdate := false

	if member.Alias != "" && logtoUser.Name == "" {
		patchRequest.Name = member.Alias
		needsUpdate = true
	}
	if member.Avatar != "" && logtoUser.Avatar == "" {
		patchRequest.Avatar = member.Avatar
		needsUpdate = true
	}

	if needsUpdate {
		_, err := service.logtoService.PatchUserById(logtoUser.Id, patchRequest)
		return err
	}
	return nil
}

var MemberServiceApp = new(MemberService)

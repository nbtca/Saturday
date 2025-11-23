package service

import (
	"net/http"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"
)

type MemberService struct{}

func (service *MemberService) GetMemberById(id string) (model.Member, error) {
	member, err := repo.GetMemberById(id)
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
	member, err := repo.GetMemberByLogtoId(logtoId)
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
	member, err := repo.GetMemberByGithubId(githubId)
	if err != nil {
		return model.Member{}, err
	}
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed").
			AddDetailError("member", "githubId", "invalid githubId")
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
	members, err := repo.GetMembers(offset, limit)
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
	members, err := repo.GetMembers(offset, limit)
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
	exist, err := repo.ExistMember(member.MemberId)
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
	if err := repo.CreateMember(member); err != nil {
		return err
	}
	return nil
}

func (service *MemberService) CreateToken(member model.Member) (string, error) {
	res, err := util.CreateToken(util.Payload{Who: member.MemberId, Member: member, Role: member.Role})
	return res, err
}

// func (service *MemberService) UpdateBasic(member model.Member) error {
// 	exist, err := repo.ExistRole(member.Role)
// 	if err != nil {
// 		return err
// 	}
// 	if !exist {
// 		return util.
// 			MakeServiceError(http.StatusUnprocessableEntity).
// 			SetMessage("Validation Failed").
// 			AddDetailError("member", "role", "invalid role")
// 	}
// 	if err := repo.UpdateMember(member); err != nil {
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
	exist, err := repo.ExistRole(member.Role)
	if err != nil {
		return err
	}
	if !exist {
		return util.
			MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("member", "role", "invalid role")
	}
	if err := repo.UpdateMember(member); err != nil {
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
	if err := repo.UpdateMember(member); err != nil {
		return err
	}
	return nil
}

// GetNotificationPreferences returns the notification preferences for a member
func (service *MemberService) GetNotificationPreferences(memberId string) ([]model.NotificationPreferenceItem, error) {
	member, err := repo.GetMemberById(memberId)
	if err != nil {
		return nil, err
	}
	if member == (model.Member{}) {
		return nil, util.MakeServiceError(http.StatusNotFound).SetMessage("Member not found")
	}

	prefs := member.GetNotificationPreferences()

	items := []model.NotificationPreferenceItem{
		{
			NotificationType: model.NotifNewEventCreated,
			Enabled:          prefs.NewEventCreated,
			Description:      model.GetNotificationDescription(model.NotifNewEventCreated),
		},
		{
			NotificationType: model.NotifEventAssignedToMe,
			Enabled:          prefs.EventAssignedToMe,
			Description:      model.GetNotificationDescription(model.NotifEventAssignedToMe),
		},
		{
			NotificationType: model.NotifEventStatusChanged,
			Enabled:          prefs.EventStatusChanged,
			Description:      model.GetNotificationDescription(model.NotifEventStatusChanged),
		},
	}

	return items, nil
}

// UpdateNotificationPreferences updates the notification preferences for a member
func (service *MemberService) UpdateNotificationPreferences(memberId string, preferences model.NotificationPreferences) error {
	exist, err := repo.ExistMember(memberId)
	if err != nil {
		return err
	}
	if !exist {
		return util.MakeServiceError(http.StatusNotFound).SetMessage("Member not found")
	}

	if err := repo.UpdateNotificationPreferences(memberId, preferences); err != nil {
		return err
	}
	return nil
}

// GetMembersWithNotificationEnabled returns members who have enabled a specific notification type
func (service *MemberService) GetMembersWithNotificationEnabled(notifType model.NotificationType) ([]model.Member, error) {
	members, err := repo.GetMembersWithNotificationEnabled(notifType)
	if err != nil {
		return nil, err
	}
	return members, nil
}

var MemberServiceApp = new(MemberService)

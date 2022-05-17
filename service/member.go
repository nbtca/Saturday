package service

import (
	"net/http"
	"saturday/model"
	"saturday/repo"
	"saturday/util"
)

type MemberService struct{}

func (service *MemberService) GetMemberById(id string) (model.Member, error) {
	member, err := repo.GetMemberById(id)
	if err != nil {
		return model.Member{}, err
	}
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
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

func (service *MemberService) CreateMember(member *model.Member) error {
	exist, err := repo.ExistMember(member.MemberId)
	if err != nil {
		return err
	}
	if exist {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
		return error
	}
	if err := repo.CreateMember(member); err != nil {
		return err
	}
	return nil
}

func (service *MemberService) CreateToken(member model.Member) (string, error) {
	res, err := util.CreateToken(util.Payload{Who: member.MemberId, Role: member.Role})
	return res, err
}

func (service *MemberService) UpdateBasic(member model.Member) error {
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

var MemberServiceApp = new(MemberService)

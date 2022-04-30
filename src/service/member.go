package service

import (
	"net/http"
	"saturday/src/model"
	"saturday/src/repo"
	"saturday/src/util"
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

func (service *MemberService) CreateMember(member *model.Member) (model.Member, error) {
	exist, err := repo.ExistMember(member.MemberId)
	if err != nil {
		return model.Member{}, err
	}
	if exist {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
		return model.Member{}, error
	}
	if err := repo.CreateMember(member); err != nil {
		return model.Member{}, err
	}
	return service.GetMemberById(member.MemberId)
}

func (service *MemberService) CreateToken(member model.Member) (string, error) {
	res := util.CreateToken(util.Payload{}) //TODO
	return res, nil
}

var MemberServiceApp = new(MemberService)

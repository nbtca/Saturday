package service

import (
	"gin-example/src/model"
	"gin-example/src/model/dto"
	"gin-example/src/repo"
	"gin-example/src/util"
	"net/http"
)

type MemberService struct {
}

func (service *MemberService) GetMemberById(id string) (model.Member, error) {
	member := repo.GetMemberById(id)
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
		return member, error
	} else {
		return member, nil
	}
}
func (service *MemberService) GetMembers(offset uint64, limit uint64) []model.Member {
	return repo.GetMembers(offset, limit)
}

func (service *MemberService) CreateMember(member *model.Member) (model.Member, error) {
	err := repo.CreateMember(member)
	if err != nil {
		return model.Member{}, err
	}
	return service.GetMemberById(member.MemberId)
}

type MemberAccount struct {
	MemberId string `json:"member_id" validate:"required,len=10,numeric"`
	Password string `json:"password" validate:"required"`
}

func (service *MemberService) CreateToken(account MemberAccount) (dto.CreateMemberTokenResponse, error) {
	// return repo.GetToken(id)
	member, err := service.GetMemberById(account.MemberId)
	if err != nil {
		return dto.CreateMemberTokenResponse{}, err
	}
	if member.Password != account.Password {
		serviceError := util.MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
		return dto.CreateMemberTokenResponse{}, serviceError
	}
	res := dto.CreateMemberTokenResponse{
		Member: member,
		Token:  "token",
	}
	return res, nil
}

var MemberServiceApp = new(MemberService)

package service

import (
	model "gin-example/models"
	"gin-example/src/repo"
	"gin-example/util"
	"net/http"
)

type MemberService struct {
	// model.Member
}

func (service *MemberService) GetMemberById(id string) (model.Member, error) {
	member := repo.GetMemberById(id)
	if member == (model.Member{}) {
		error := util.MakeServiceError(http.StatusUnprocessableEntity)
		return member, error
	} else {
		return member, nil
	}
}

func (service *MemberService) GetMembers(offset uint64, limit uint64) []model.Member {
	return repo.GetMembers(offset, limit)
}

var MemberServiceApp = new(MemberService)

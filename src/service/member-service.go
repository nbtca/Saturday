package service

import (
	model "gin-example/models"
	"gin-example/src/repo"
	"gin-example/util"
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

func (service *MemberService) GetToken(id string) (string, error) {
	// return repo.GetToken(id)
	return "TOKEN", nil
}

var MemberServiceApp = new(MemberService)

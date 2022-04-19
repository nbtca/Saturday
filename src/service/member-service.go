package service

import (
	model "gin-example/models"
	"gin-example/src/repo"
)

type MemberService struct {
	Repo *repo.MemberRepo
}

func (service *MemberService) GetMemberById(id string) model.Member {
	return service.Repo.GetMemberById(id)
}

func (MemberService *MemberService) GetMembers(offset uint64, limit uint64) []model.Member {
	return MemberService.Repo.GetMembers(offset, limit)
}

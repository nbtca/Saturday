package middleware

import (
	"github.com/nbtca/saturday/service"
)

type Role string

const (
	Member Role = "member"
	Admin  Role = "admin"
)

type AuthContextUser struct {
	UserInfo service.FetchUserInfoResponse
	Role     []string
}

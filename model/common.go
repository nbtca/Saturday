package model

type Identity struct {
	Id       string
	ClientId int64
	Member   Member
	Role     string
}

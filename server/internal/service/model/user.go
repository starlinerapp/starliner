package model

import "starliner.app/internal/domain"

type User struct {
	Id           int64
	BetterAuthId string
}

func NewUser(u *domain.User) *User {
	return &User{
		Id:           u.Id,
		BetterAuthId: u.BetterAuthId,
	}
}

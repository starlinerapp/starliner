package value

import (
	"starliner.app/internal/api/domain/entity"
)

type User struct {
	Id           int64
	BetterAuthId string
}

func NewUser(u *entity.User) *User {
	return &User{
		Id:           u.Id,
		BetterAuthId: u.BetterAuthId,
	}
}

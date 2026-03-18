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

func NewUsers(us []*entity.User) []*User {
	users := make([]*User, len(us))
	for i, u := range us {
		users[i] = NewUser(u)
	}
	return users
}

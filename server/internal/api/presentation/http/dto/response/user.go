package response

import "starliner.app/internal/api/domain/value"

type User struct {
	UserId       int64  `json:"user_id" binding:"required"`
	BetterAuthId string `json:"better_auth_id" binding:"required"`
}

func NewUser(user *value.User) User {
	return User{
		UserId:       user.Id,
		BetterAuthId: user.BetterAuthId,
	}
}

func NewUsers(users []*value.User) []User {
	result := make([]User, len(users))
	for i, user := range users {
		result[i] = NewUser(user)
	}
	return result
}

package response

type User struct {
	UserId       int64  `json:"user_id" binding:"required"`
	BetterAuthId string `json:"better_auth_id" binding:"required"`
}

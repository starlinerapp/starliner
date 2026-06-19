package request

type RegisterRunner struct {
	Token string `json:"token" binding:"required"`
}

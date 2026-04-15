package request

type CreateTeam struct {
	Slug string `json:"slug" binding:"required,max=50"`
}

type JoinTeam struct {
	Slug string `json:"slug" binding:"required"`
}

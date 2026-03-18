package request

type CreateTeam struct {
	Name string `json:"name" binding:"required"`
}

type JoinTeam struct {
	Slug string `json:"slug" binding:"required"`
}

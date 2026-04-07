package request

type CreateTeam struct {
	Slug string `json:"slug" binding:"required,max=50"`
}

type JoinTeam struct {
	Slug string `json:"slug" binding:"required"`
}

type AddTeamMember struct {
	UserID int64 `json:"userId" binding:"required"`
}

type RemoveTeamMember struct {
	UserID int64 `json:"userId" binding:"required"`
}

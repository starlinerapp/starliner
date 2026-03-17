package response

import "starliner.app/internal/api/domain/value"

type Team struct {
	Id             int64  `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Slug           string `json:"slug" binding:"required"`
	OrganizationId int64  `json:"organization_id" binding:"required"`
}

func NewTeam(team *value.Team) *Team {
	return &Team{
		Id:             team.Id,
		Name:           team.Name,
		Slug:           team.Slug,
		OrganizationId: team.OrganizationId,
	}
}

func NewTeams(teams []*value.Team) []*Team {
	res := make([]*Team, len(teams))
	for i, team := range teams {
		res[i] = NewTeam(team)
	}

	return res
}

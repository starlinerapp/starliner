package value

import "starliner.app/internal/api/domain/entity"

type Team struct {
	Id             int64
	Slug           string
	OrganizationId int64
}

func NewTeam(t *entity.Team) *Team {
	return &Team{
		Id:             t.Id,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationId,
	}
}

func NewTeams(ts []*entity.Team) []*Team {
	teams := make([]*Team, len(ts))
	for i, t := range ts {
		teams[i] = NewTeam(t)
	}
	return teams
}

package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, name string, slug string, organizationID int64) (*entity.Team, error)
	DeleteTeam(ctx context.Context, id int64) error
	GetTeamBySlug(ctx context.Context, slug string, organizationID int64) (*entity.Team, error)
	GetUserTeams(ctx context.Context, organizationID int64, userID int64) ([]*entity.Team, error)
	GetTeamMembers(ctx context.Context, teamID int64) ([]*entity.User, error)
	AddTeamMember(ctx context.Context, teamID int64, userID int64) error
	RemoveTeamMember(ctx context.Context, teamID int64, userID int64) error
}

package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, slug string, organizationID int64) (*entity.Team, error)
	DeleteTeam(ctx context.Context, id int64) error
	DeleteTeamIfEmpty(ctx context.Context, id int64) error
	GetTeamBySlug(ctx context.Context, slug string, organizationID int64) (*entity.Team, error)
	GetUserTeams(ctx context.Context, organizationID int64, userID int64) ([]*entity.Team, error)
	GetTeamMembers(ctx context.Context, teamID int64) ([]*entity.User, error)
	AddTeamMember(ctx context.Context, teamID int64, userID int64) error
	RemoveTeamMember(ctx context.Context, teamID int64, userID int64) error
	GetTeamById(ctx context.Context, id int64) (*entity.Team, error)
	FindTeamByIdAndUserId(ctx context.Context, teamID int64, userID int64) (*entity.Team, error)
	SetTeamRepositories(ctx context.Context, teamID int64, repos []*value.TeamRepo, githubAppID int64) error
	GetTeamRepositories(ctx context.Context, teamID int64) ([]*entity.TeamRepository, error)
	GetTeamClusters(ctx context.Context, teamID int64) ([]*entity.TeamCluster, error)
	SetTeamClusters(ctx context.Context, teamID int64, clusters []*value.TeamCluster) error
	AssignClusterToTeam(ctx context.Context, teamID int64, clusterId int64) error
	GetTeamCluster(ctx context.Context, teamID int64, clusterId int64) (*entity.TeamCluster, error)
}

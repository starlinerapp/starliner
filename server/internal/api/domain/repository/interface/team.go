package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
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
	AssignRepoToTeam(ctx context.Context, teamID int64, githubRepoID int64, repoName string, githubAppID int64) error
	UnassignRepoFromTeam(ctx context.Context, teamID int64, githubRepoID int64) error
	GetTeamRepositories(ctx context.Context, teamID int64) ([]*entity.TeamRepository, error)
	GetTeamClusters(ctx context.Context, teamID int64) ([]*entity.TeamCluster, error)
	AssignClusterToTeam(ctx context.Context, teamID int64, clusterId int64) error
	UnassignClusterFromTeam(ctx context.Context, teamID int64, clusterId int64) error
	GetTeamCluster(ctx context.Context, teamID int64, clusterId int64) (*entity.TeamCluster, error)
}

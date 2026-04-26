package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
)

type TeamRepository struct {
	queries *sqlc.Queries
}

func (tr *TeamRepository) GetTeamById(ctx context.Context, id int64) (*entity.Team, error) {
	t, err := tr.queries.GetTeamById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr *TeamRepository) CreateTeam(ctx context.Context, slug string, organizationID int64) (*entity.Team, error) {
	t, err := tr.queries.CreateTeam(ctx, sqlc.CreateTeamParams{
		Slug:           slug,
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr *TeamRepository) DeleteTeam(ctx context.Context, id int64) error {
	err := tr.queries.DeleteTeam(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TeamRepository) GetTeamBySlug(ctx context.Context, slug string, organizationID int64) (*entity.Team, error) {
	t, err := tr.queries.GetTeamBySlug(ctx, sqlc.GetTeamBySlugParams{
		Slug:           slug,
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr *TeamRepository) GetUserTeams(ctx context.Context, organizationID int64, userID int64) ([]*entity.Team, error) {
	rows, err := tr.queries.GetUserTeams(ctx, sqlc.GetUserTeamsParams{
		OrganizationID: organizationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}

	teams := make([]*entity.Team, len(rows))
	for i, row := range rows {
		teams[i] = &entity.Team{
			Id:             row.ID,
			Slug:           row.Slug,
			OrganizationId: row.OrganizationID,
		}
	}

	return teams, nil
}

func (tr *TeamRepository) GetTeamMembers(ctx context.Context, teamID int64) ([]*entity.User, error) {
	rows, err := tr.queries.GetTeamMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(rows))
	for i, row := range rows {
		users[i] = &entity.User{
			Id:           row.ID,
			BetterAuthId: row.BetterAuthID,
		}
	}

	return users, nil
}

func (tr *TeamRepository) AddTeamMember(ctx context.Context, teamID int64, userID int64) error {
	err := tr.queries.AddTeamMember(ctx, sqlc.AddTeamMemberParams{
		TeamID: teamID,
		UserID: userID,
	})

	return err
}

func (tr *TeamRepository) RemoveTeamMember(ctx context.Context, teamID int64, userID int64) error {
	err := tr.queries.RemoveTeamMember(ctx, sqlc.RemoveTeamMemberParams{
		TeamID: teamID,
		UserID: userID,
	})

	return err
}

func (tr *TeamRepository) FindTeamByIdAndUserId(ctx context.Context, teamID int64, userID int64) (*entity.Team, error) {
	t, err := tr.queries.FindTeamByIdAndUserId(ctx, sqlc.FindTeamByIdAndUserIdParams{
		ID:     teamID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr *TeamRepository) DeleteTeamIfEmpty(ctx context.Context, id int64) error {
	return tr.queries.DeleteTeamIfEmpty(ctx, id)
}

func (tr *TeamRepository) AssignRepoToTeam(ctx context.Context, teamID int64, githubRepoID int64, repoName string, githubAppID int64) error {
	return tr.queries.AssignRepoToTeam(ctx, sqlc.AssignRepoToTeamParams{
		TeamID:       teamID,
		GithubRepoID: githubRepoID,
		RepoName:     repoName,
		GithubAppID:  githubAppID,
	})
}

func (tr *TeamRepository) UnassignRepoFromTeam(ctx context.Context, teamID int64, githubRepoID int64) error {
	return tr.queries.UnassignRepoFromTeam(ctx, sqlc.UnassignRepoFromTeamParams{
		TeamID:       teamID,
		GithubRepoID: githubRepoID,
	})
}

func (tr *TeamRepository) GetTeamRepositories(ctx context.Context, teamID int64) ([]*entity.TeamRepository, error) {
	rows, err := tr.queries.GetTeamRepositories(ctx, teamID)
	if err != nil {
		return nil, err
	}

	repos := make([]*entity.TeamRepository, len(rows))
	for i, row := range rows {
		repos[i] = &entity.TeamRepository{
			TeamId:       teamID,
			GithubRepoId: row.GithubRepoID,
			RepoName:     row.RepoName,
		}
	}

	return repos, nil
}

func (tr *TeamRepository) GetTeamClusters(ctx context.Context, teamID int64) ([]*entity.TeamCluster, error) {
	rows, err := tr.queries.GetTeamClusters(ctx, teamID)
	if err != nil {
		return nil, err
	}

	clusters := make([]*entity.TeamCluster, len(rows))
	for i, row := range rows {
		clusters[i] = &entity.TeamCluster{
			TeamId:      teamID,
			ClusterId:   row.ID,
			ClusterName: row.Name,
			ServerType:  row.ServerType,
		}
	}
	return clusters, nil
}

func (tr *TeamRepository) AssignClusterToTeam(ctx context.Context, teamID int64, clusterId int64) error {
	err := tr.queries.AssignTeamCluster(ctx, sqlc.AssignTeamClusterParams{
		TeamID:    teamID,
		ClusterID: clusterId,
	})
	return err
}

func (tr *TeamRepository) UnassignClusterFromTeam(ctx context.Context, teamID int64, clusterId int64) error {
	err := tr.queries.UnassignTeamCluster(ctx, sqlc.UnassignTeamClusterParams{
		TeamID:    teamID,
		ClusterID: clusterId,
	})
	return err
}

func (tr *TeamRepository) GetTeamCluster(ctx context.Context, teamID int64, clusterId int64) (*entity.TeamCluster, error) {
	cluster, err := tr.queries.GetTeamCluster(ctx, sqlc.GetTeamClusterParams{
		TeamID:    teamID,
		ClusterID: clusterId,
	})
	if err != nil {
		return nil, err
	}
	return &entity.TeamCluster{
		TeamId:      teamID,
		ClusterId:   cluster.ID,
		ClusterName: cluster.Name,
		ServerType:  cluster.ServerType,
	}, nil
}

var _ interfaces.TeamRepository = (*TeamRepository)(nil)

func NewTeamRepository(queries *sqlc.Queries) interfaces.TeamRepository {
	return &TeamRepository{queries: queries}
}

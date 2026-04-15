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

func (tr TeamRepository) GetTeamById(ctx context.Context, id int64) (*entity.Team, error) {
	t, err := tr.queries.GetTeamById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Name:           t.Name,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr TeamRepository) CreateTeam(ctx context.Context, name string, slug string, organizationID int64) (*entity.Team, error) {
	t, err := tr.queries.CreateTeam(ctx, sqlc.CreateTeamParams{
		Name:           name,
		Slug:           slug,
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Name:           t.Name,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr TeamRepository) DeleteTeam(ctx context.Context, id int64) error {
	err := tr.queries.DeleteTeam(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (tr TeamRepository) GetTeamBySlug(ctx context.Context, slug string, organizationID int64) (*entity.Team, error) {
	t, err := tr.queries.GetTeamBySlug(ctx, sqlc.GetTeamBySlugParams{
		Slug:           slug,
		OrganizationID: organizationID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Name:           t.Name,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr TeamRepository) GetUserTeams(ctx context.Context, organizationID int64, userID int64) ([]*entity.Team, error) {
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
			Name:           row.Name,
			Slug:           row.Slug,
			OrganizationId: row.OrganizationID,
		}
	}

	return teams, nil
}

func (tr TeamRepository) GetTeamMembers(ctx context.Context, teamID int64) ([]*entity.User, error) {
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

func (tr TeamRepository) AddTeamMember(ctx context.Context, teamID int64, userID int64) error {
	err := tr.queries.AddTeamMember(ctx, sqlc.AddTeamMemberParams{
		TeamID: teamID,
		UserID: userID,
	})

	return err
}

func (tr TeamRepository) RemoveTeamMember(ctx context.Context, teamID int64, userID int64) error {
	err := tr.queries.RemoveTeamMember(ctx, sqlc.RemoveTeamMemberParams{
		TeamID: teamID,
		UserID: userID,
	})

	return err
}

func (tr TeamRepository) FindTeamByIdAndUserId(ctx context.Context, teamID int64, userID int64) (*entity.Team, error) {
	t, err := tr.queries.FindTeamByIdAndUserId(ctx, sqlc.FindTeamByIdAndUserIdParams{
		ID:     teamID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Team{
		Id:             t.ID,
		Name:           t.Name,
		Slug:           t.Slug,
		OrganizationId: t.OrganizationID,
	}, nil
}

func (tr TeamRepository) DeleteTeamIfEmpty(ctx context.Context, id int64) error {
	return tr.queries.DeleteTeamIfEmpty(ctx, id)
}

func (tr TeamRepository) AssignRepoToTeam(ctx context.Context, teamID int64, githubRepoID int64, repoName string) error {
	return tr.queries.AssignRepoToTeam(ctx, sqlc.AssignRepoToTeamParams{
		TeamID:       teamID,
		GithubRepoID: githubRepoID,
		RepoName:     repoName,
	})
}

func (tr TeamRepository) UnassignRepoFromTeam(ctx context.Context, teamID int64, githubRepoID int64) error {
	return tr.queries.UnassignRepoFromTeam(ctx, sqlc.UnassignRepoFromTeamParams{
		TeamID:       teamID,
		GithubRepoID: githubRepoID,
	})
}

func (tr TeamRepository) GetTeamRepositories(ctx context.Context, teamID int64) ([]*entity.TeamRepository, error) {
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

var _ interfaces.TeamRepository = (*TeamRepository)(nil)

func NewTeamRepository(queries *sqlc.Queries) interfaces.TeamRepository {
	return &TeamRepository{queries: queries}
}

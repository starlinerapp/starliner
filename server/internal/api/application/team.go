package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	coreService "starliner.app/internal/core/domain/service"
)

type TeamApplication struct {
	teamRepository       interfaces.TeamRepository
	organizationService  *service.OrganizationService
	normalizationService *coreService.NormalizerService
}

func NewTeamApplication(
	teamRepository interfaces.TeamRepository,
	organizationService *service.OrganizationService,
	normalizationService *coreService.NormalizerService,
) *TeamApplication {
	return &TeamApplication{
		teamRepository:       teamRepository,
		organizationService:  organizationService,
		normalizationService: normalizationService,
	}
}

func (ts *TeamApplication) CreateTeam(ctx context.Context, name string, organizationId int64, userId int64) (*value.Team, error) {
	err := ts.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	teamSlug, err := ts.normalizationService.FormatToDNS1123(name)
	if err != nil {
		return nil, err
	}

	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		return nil, err
	}
	teamSlug = fmt.Sprintf("%s-%s", teamSlug, hex.EncodeToString(suffix))

	team, err := ts.teamRepository.CreateTeam(
		ctx,
		name,
		teamSlug,
		organizationId,
	)
	if err != nil {
		return nil, err
	}

	err = ts.teamRepository.AddTeamMember(ctx, team.Id, userId)
	if err != nil {
		return nil, err
	}

	return value.NewTeam(team), nil
}

func (ts *TeamApplication) GetUserTeams(ctx context.Context, organizationId int64, userId int64) ([]*value.Team, error) {
	err := ts.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	teams, err := ts.teamRepository.GetUserTeams(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	return value.NewTeams(teams), nil
}

func (ts *TeamApplication) GetTeamMembers(ctx context.Context, userId int64, teamId int64) ([]*value.User, error) {
	_, err := ts.teamRepository.FindTeamByIdAndUserId(ctx, teamId, userId)
	if err != nil {
		return nil, err
	}

	members, err := ts.teamRepository.GetTeamMembers(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return value.NewUsers(members), nil
}

func (ts *TeamApplication) AddTeamMember(ctx context.Context, userId int64, teamId int64) error {
	_, err := ts.teamRepository.FindTeamByIdAndUserId(ctx, teamId, userId)
	if err != nil {
		return err
	}

	return ts.teamRepository.AddTeamMember(ctx, teamId, userId)
}

func (ts *TeamApplication) JoinTeam(ctx context.Context, organizationId int64, userId int64, teamSlug string) error {
	err := ts.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	team, err := ts.teamRepository.GetTeamBySlug(ctx, teamSlug, organizationId)
	if err != nil {
		return err
	}

	return ts.teamRepository.AddTeamMember(ctx, team.Id, userId)
}

func (ts *TeamApplication) RemoveTeamMember(ctx context.Context, userId int64, teamId int64) error {
	_, err := ts.teamRepository.FindTeamByIdAndUserId(ctx, teamId, userId)
	if err != nil {
		return err
	}

	err = ts.teamRepository.RemoveTeamMember(ctx, teamId, userId)
	if err != nil {
		return err
	}

	return ts.teamRepository.DeleteTeamIfEmpty(ctx, teamId)
}

func (ts *TeamApplication) AssignRepoToTeam(ctx context.Context, organizationId int64, userId int64, teamId int64, githubRepoId int64, repoName string) error {
	err := ts.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	return ts.teamRepository.AssignRepoToTeam(ctx, teamId, githubRepoId, repoName)
}

func (ts *TeamApplication) UnassignRepoFromTeam(ctx context.Context, organizationId int64, userId int64, teamId int64, githubRepoId int64) error {
	err := ts.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	return ts.teamRepository.UnassignRepoFromTeam(ctx, teamId, githubRepoId)
}

func (ts *TeamApplication) GetTeamRepositories(ctx context.Context, organizationId int64, userId int64, teamId int64) ([]*value.TeamRepo, error) {
	err := ts.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	repos, err := ts.teamRepository.GetTeamRepositories(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return value.NewTeamRepos(repos), nil
}

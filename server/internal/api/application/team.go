package application

import (
	"context"
	"errors"

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

func (ta *TeamApplication) CreateTeam(ctx context.Context, slug string, organizationId int64, userId int64) (*value.Team, error) {
	err := ta.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	normalized, err := ta.normalizationService.FormatToDNS1123(slug)
	if err != nil {
		return nil, err
	}

	if normalized != slug {
		return nil, errors.New("invalid slug format")
	}

	team, err := ta.teamRepository.CreateTeam(
		ctx,
		slug,
		organizationId,
	)
	if err != nil {
		return nil, err
	}

	err = ta.teamRepository.AddTeamMember(ctx, team.Id, userId)
	if err != nil {
		return nil, err
	}

	return value.NewTeam(team), nil
}

func (ta *TeamApplication) GetUserTeams(ctx context.Context, organizationId int64, userId int64) ([]*value.Team, error) {
	err := ta.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	teams, err := ta.teamRepository.GetUserTeams(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	return value.NewTeams(teams), nil
}

func (ta *TeamApplication) GetTeamMembers(ctx context.Context, userId int64, teamId int64) ([]*value.User, error) {
	_, err := ta.teamRepository.FindTeamByIdAndUserId(ctx, teamId, userId)
	if err != nil {
		return nil, err
	}

	members, err := ta.teamRepository.GetTeamMembers(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return value.NewUsers(members), nil
}

func (ta *TeamApplication) AddTeamMember(ctx context.Context, userId int64, teamId int64, callerId int64) error {
	team, err := ta.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return err
	}

	err = ta.organizationService.ValidateUserOrgOwner(ctx, team.OrganizationId, callerId)
	if err != nil {
		return err
	}

	return ta.teamRepository.AddTeamMember(ctx, teamId, userId)
}

func (ta *TeamApplication) JoinTeam(ctx context.Context, organizationId int64, userId int64, teamSlug string) error {
	err := ta.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	team, err := ta.teamRepository.GetTeamBySlug(ctx, teamSlug, organizationId)
	if err != nil {
		return err
	}

	return ta.teamRepository.AddTeamMember(ctx, team.Id, userId)
}

func (ta *TeamApplication) RemoveTeamMember(ctx context.Context, userId int64, teamId int64, callerId int64) error {
	team, err := ta.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return err
	}

	err = ta.organizationService.ValidateUserOrgOwner(ctx, team.OrganizationId, callerId)
	if err != nil {
		return err
	}

	// Only the org owner is allowed to manage team members and is part of every team.
	// This check enforces the owner does not remove himself from a team.
	if callerId == userId {
		return errors.New("organization owner cannot be removed from team")
	}

	err = ta.teamRepository.RemoveTeamMember(ctx, teamId, userId)
	if err != nil {
		return err
	}

	return ta.teamRepository.DeleteTeamIfEmpty(ctx, teamId)
}

func (ta *TeamApplication) AssignRepoToTeam(ctx context.Context, userId int64, teamId int64, githubRepoId int64, repoName string) error {
	team, err := ta.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return err
	}

	err = ta.organizationService.ValidateUserOrgOwner(ctx, team.OrganizationId, userId)
	if err != nil {
		return err
	}

	return ta.teamRepository.AssignRepoToTeam(ctx, teamId, githubRepoId, repoName)
}

func (ta *TeamApplication) UnassignRepoFromTeam(ctx context.Context, userId int64, teamId int64, githubRepoId int64) error {
	team, err := ta.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return err
	}

	err = ta.organizationService.ValidateUserOrgOwner(ctx, team.OrganizationId, userId)
	if err != nil {
		return err
	}

	return ta.teamRepository.UnassignRepoFromTeam(ctx, teamId, githubRepoId)
}

func (ta *TeamApplication) GetTeamRepositories(ctx context.Context, userId int64, teamId int64) ([]*value.TeamRepo, error) {
	_, err := ta.teamRepository.FindTeamByIdAndUserId(ctx, teamId, userId)
	if err != nil {
		return nil, err
	}
	repos, err := ta.teamRepository.GetTeamRepositories(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return value.NewTeamRepos(repos), nil
}

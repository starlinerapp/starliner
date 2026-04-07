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

func (ts *TeamApplication) CreateTeam(ctx context.Context, slug string, organizationId int64, userId int64) (*value.Team, error) {
	err := ts.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	normalized, err := ts.normalizationService.FormatToDNS1123(slug)
	if err != nil {
		return nil, err
	}

	if normalized != slug {
		return nil, errors.New("invalid slug format")
	}

	team, err := ts.teamRepository.CreateTeam(
		ctx,
		slug,
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
	err := ts.teamRepository.ValidateUserTeamAccess(ctx, teamId, userId)
	if err != nil {
		return nil, err
	}

	members, err := ts.teamRepository.GetTeamMembers(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return value.NewUsers(members), nil
}

func (ts *TeamApplication) AddTeamMember(ctx context.Context, userId int64, teamId int64, callerId int64) error {
	team, err := ts.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return err
	}

	err = ts.organizationService.ValidateUserOrgOwner(ctx, team.OrganizationId, callerId)
	if err != nil {
		return err
	}

	err = ts.organizationService.ValidateUserInOrg(ctx, team.OrganizationId, userId)
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

func (ts *TeamApplication) RemoveTeamMember(ctx context.Context, userId int64, teamId int64, callerId int64) error {
	team, err := ts.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return err
	}

	err = ts.organizationService.ValidateUserOrgOwner(ctx, team.OrganizationId, callerId)
	if err != nil {
		return err
	}

	// Only the org owner is allowed to manage team members and is part of every team.
	// This check enforces the owner does not remove himself from a team.
	if callerId == userId {
		return errors.New("organization owner cannot be removed from team")
	}

	err = ts.teamRepository.RemoveTeamMember(ctx, teamId, userId)
	if err != nil {
		return err
	}

	return ts.teamRepository.DeleteTeamIfEmpty(ctx, teamId)
}

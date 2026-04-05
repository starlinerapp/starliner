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

func (ts *TeamApplication) AddTeamMember(ctx context.Context, userId int64, teamId int64) error {
	err := ts.teamRepository.ValidateUserTeamAccess(ctx, teamId, userId)
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
	err := ts.teamRepository.ValidateUserTeamAccess(ctx, teamId, userId)
	if err != nil {
		return err
	}

	err = ts.teamRepository.RemoveTeamMember(ctx, teamId, userId)
	if err != nil {
		return err
	}

	return ts.teamRepository.DeleteTeamIfEmpty(ctx, teamId)
}

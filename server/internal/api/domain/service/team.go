package service

import (
	"context"
	_interface "starliner.app/internal/api/domain/repository/interface"
)

type TeamService struct {
	clusterRepository _interface.ClusterRepository
	teamRepository    _interface.TeamRepository
}

func NewTeamService(
	clusterRepository _interface.ClusterRepository,
	teamRepository _interface.TeamRepository,
) *TeamService {
	return &TeamService{
		clusterRepository: clusterRepository,
		teamRepository:    teamRepository,
	}
}

func (os *TeamService) ValidateUserInTeam(ctx context.Context, userId int64, teamId int64) error {
	_, err := os.teamRepository.FindTeamByIdAndUserId(ctx, teamId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (os *TeamService) ValidateUserAndClusterInTeam(ctx context.Context, userId int64, teamId int64, clusterId int64) error {
	err := os.ValidateUserInTeam(ctx, userId, teamId)
	if err != nil {
		return err
	}
	_, err = os.teamRepository.FindTeamCluster(ctx, teamId, clusterId)
	if err != nil {
		return err
	}
	return nil
}

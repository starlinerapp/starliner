package service

import (
	"context"

	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
)

type BuildService struct {
	buildRepository interfaces.BuildRepository
}

func NewBuildService(buildRepository interfaces.BuildRepository) *BuildService {
	return &BuildService{buildRepository: buildRepository}
}

func (s *BuildService) CreateDeployOnlyBuild(
	ctx context.Context,
	deploymentId int64,
	source string,
) error {
	b, err := s.buildRepository.CreateBuild(ctx, deploymentId, source)
	if err != nil {
		return err
	}

	return s.buildRepository.UpdateBuild(
		ctx,
		b.Id,
		value.BuildStatusSuccess,
		nil,
		nil,
		"Build skipped",
	)
}

package application

import (
	"context"
	interfaces "starliner.app/internal/api/domain/repository/interface"
)

type BuildApplication struct {
	buildRepository interfaces.BuildRepository
}

func NewBuildApplication(buildRepository interfaces.BuildRepository) *BuildApplication {
	return &BuildApplication{
		buildRepository: buildRepository,
	}
}

func (ba *BuildApplication) GetBuildLogs(ctx context.Context, userId int64, buildId int64) (*string, error) {
	return ba.buildRepository.GetBuildLogs(ctx, userId, buildId)
}

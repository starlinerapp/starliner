package service

import (
	"context"
	"errors"
	interfaces "starliner.app/internal/domain/repository/interface"
)

type ProjectService struct {
	projectRepository interfaces.ProjectRepository
}

func NewProjectService(projectRepository interfaces.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepository: projectRepository}
}

func (ps *ProjectService) ValidateUserHasPermission(ctx context.Context, projectId int64, userId int64) error {
	project, err := ps.projectRepository.GetUserProject(ctx, userId, projectId)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("user does not have permission")
	}
	return nil
}

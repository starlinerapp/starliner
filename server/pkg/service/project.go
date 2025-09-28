package service

import (
	"context"
	"errors"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/repository"
)

type ProjectService struct {
	projectRepository      *repository.ProjectRepository
	organizationRepository *repository.OrganizationRepository
}

func NewProjectService(projectRepository *repository.ProjectRepository, organizationRepository *repository.OrganizationRepository) *ProjectService {
	return &ProjectService{projectRepository: projectRepository, organizationRepository: organizationRepository}
}

func (ps *ProjectService) CreateProject(ctx context.Context, name string, organizationId int64, userId int64) (*domain.Project, error) {
	organizations, err := ps.organizationRepository.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, err
	}

	found := false
	for _, org := range organizations {
		if org.Id == organizationId {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("organization not found")
	}

	return ps.projectRepository.CreateProject(ctx, name, organizationId)
}

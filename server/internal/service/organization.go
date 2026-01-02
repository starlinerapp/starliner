package service

import (
	"context"
	"errors"
	interfaces "starliner.app/internal/repository/interface"
	"starliner.app/internal/service/model"
	"strings"
)

type OrganizationService struct {
	organizationRepository interfaces.OrganizationRepository
}

func NewOrganizationService(organizationRepository interfaces.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		organizationRepository: organizationRepository,
	}
}

func (os *OrganizationService) CreateOrganization(ctx context.Context, name string, ownerID int64) error {
	trimmed := strings.TrimSpace(name)
	organizationSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	_, err := os.organizationRepository.CreateOrganization(ctx, name, organizationSlug, ownerID)
	if err != nil {
		return err
	}
	return nil
}

func (os *OrganizationService) GetUserOrganizations(ctx context.Context, userID int64) ([]model.Organization, error) {
	organizations, err := os.organizationRepository.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}

	organizationModels := make([]model.Organization, len(organizations))
	for i, org := range organizations {
		organizationModels[i] = model.Organization{
			Id:      org.Id,
			Name:    org.Name,
			Slug:    org.Slug,
			OwnerId: org.OwnerId,
		}
	}
	return organizationModels, nil
}

func (os *OrganizationService) GetProjectsForUser(ctx context.Context, userID int64, organizationID int64) ([]model.Project, error) {
	err := os.ValidateUserOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	projects, err := os.organizationRepository.GetOrganizationProjects(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	projectsModels := make([]model.Project, len(projects))
	for i, p := range projects {
		environmentModel := make([]model.Environment, len(p.Environments))
		for j, env := range p.Environments {
			environmentModel[j] = model.Environment{
				Id:   env.Id,
				Slug: env.Slug,
				Name: env.Name,
			}
		}

		projectsModels[i] = model.Project{
			Id:             p.Id,
			Name:           p.Name,
			Environments:   environmentModel,
			OrganizationId: p.OrganizationId,
		}
	}

	return projectsModels, nil
}

func (os *OrganizationService) GetClustersForUser(ctx context.Context, userID int64, organizationID int64) ([]model.Cluster, error) {
	err := os.ValidateUserOrganization(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	clusters, err := os.organizationRepository.GetOrganizationClusters(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	clusterModels := make([]model.Cluster, len(clusters))
	for i, org := range clusters {
		clusterModels[i] = model.Cluster{
			Id:             org.Id,
			Name:           org.Name,
			Status:         org.Status,
			IPv4Address:    org.IPv4Address,
			OrganizationId: org.OrganizationId,
		}
	}
	return clusterModels, nil
}

func (os *OrganizationService) ValidateUserOrganization(ctx context.Context, organizationId int64, userId int64) error {
	organizations, err := os.organizationRepository.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil
	}

	found := false
	for _, org := range organizations {
		if org.Id == organizationId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("organization not found")
	}

	return nil
}

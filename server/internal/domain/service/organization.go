package service

import (
	"context"
	"errors"
	interfaces "starliner.app/internal/domain/repository/interface"
)

type OrganizationService struct {
	organizationRepository interfaces.OrganizationRepository
}

func NewOrganizationService(
	organizationRepository interfaces.OrganizationRepository,
) *OrganizationService {
	return &OrganizationService{
		organizationRepository: organizationRepository,
	}
}

func (os *OrganizationService) ValidateUserInOrg(ctx context.Context, organizationId int64, userId int64) error {
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

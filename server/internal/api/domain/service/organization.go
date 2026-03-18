package service

import (
	"context"
	"errors"
	_interface "starliner.app/internal/api/domain/repository/interface"
)

type OrganizationService struct {
	organizationRepository _interface.OrganizationRepository
}

func NewOrganizationService(
	organizationRepository _interface.OrganizationRepository,
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
		return errors.New("user not in organization")
	}
	return nil
}

func (os *OrganizationService) ValidateUserOrgOwner(ctx context.Context, organizationId int64, userId int64) error {
	organizations, err := os.organizationRepository.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil
	}
	found := false
	for _, org := range organizations {
		if org.Id == organizationId && org.OwnerId == userId {
			found = true
		}
	}

	if !found {
		return errors.New("user does not own organization")
	}
	return nil
}

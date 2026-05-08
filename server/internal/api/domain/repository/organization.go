package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type OrganizationRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.OrganizationRepository = (*OrganizationRepository)(nil)

func NewOrganizationRepository(queries *sqlc.Queries) interfaces.OrganizationRepository {
	return &OrganizationRepository{queries: queries}
}

func (or *OrganizationRepository) CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*entity.Organization, error) {
	organization, err := or.queries.CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		Name:    name,
		Slug:    slug,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Organization{
		Id:      organization.ID,
		Name:    organization.Name,
		Slug:    organization.Slug,
		OwnerId: organization.OwnerID,
	}, nil
}

func (or *OrganizationRepository) GetOrganization(ctx context.Context, id int64) (*entity.Organization, error) {
	organization, err := or.queries.GetOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Organization{
		Id:      organization.ID,
		Name:    organization.Name,
		Slug:    organization.Slug,
		OwnerId: organization.OwnerID,
	}, nil
}

func (or *OrganizationRepository) GetUserOrganizations(ctx context.Context, userId int64) ([]*entity.Organization, error) {
	organizations, err := or.queries.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Organization, len(organizations))
	for i, organization := range organizations {
		result[i] = &entity.Organization{
			Id:      organization.ID,
			Name:    organization.Name,
			Slug:    organization.Slug,
			OwnerId: organization.OwnerID,
		}
	}

	return result, nil
}

func (or *OrganizationRepository) GetUserProjects(ctx context.Context, organizationID int64, userID int64) ([]*entity.Project, error) {
	rows, err := or.queries.GetUserProjects(ctx, sqlc.GetUserProjectsParams{
		OrganizationID: organizationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}

	projectsMap := make(map[int64]*entity.Project)
	for _, row := range rows {
		proj, exists := projectsMap[row.ID]
		if !exists {
			proj = &entity.Project{
				Id:           row.ID,
				Name:         row.Name,
				Environments: []*entity.Environment{},
				TeamId:       row.TeamID,
				TeamSlug:     row.TeamSlug,
				CreatedAt:    row.CreatedAt,
			}
			projectsMap[proj.Id] = proj
		}

		proj.Environments = append(proj.Environments, &entity.Environment{
			Id:   row.EnvironmentID,
			Slug: row.EnvironmentSlug,
			Name: row.EnvironmentName,
		})
	}

	projects := make([]*entity.Project, 0, len(projectsMap))
	for _, p := range projectsMap {
		projects = append(projects, p)
	}

	return projects, nil
}

func (or *OrganizationRepository) GetOrganizationClusters(ctx context.Context, organizationID int64) ([]*entity.Cluster, error) {
	rows, err := or.queries.GetOrganizationClusters(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return []*entity.Cluster{}, nil
	}

	clusters := make([]*entity.Cluster, 0, len(rows))
	for _, c := range rows {
		clusters = append(clusters, &entity.Cluster{
			Id:             c.ID,
			Name:           c.Name,
			TeamSlug:       utils.PtrFromNullString(c.TeamSlug),
			OrganizationId: c.OrganizationID,
			ServerType:     entity.ServerType(c.ServerType),
			CreatedAt:      c.CreatedAt,
		})
	}

	return clusters, nil
}

func (or *OrganizationRepository) UpsertProvisioningCredentials(
	ctx context.Context,
	organizationID int64,
	apiKey string,
	provider value.ProvisioningCredentialProvider,
) error {
	err := or.queries.UpsertProvisioningCredential(ctx, sqlc.UpsertProvisioningCredentialParams{
		OrganizationID: organizationID,
		Provider:       sqlc.Provider(provider),
		Secret:         apiKey,
	})

	return err
}

func (or *OrganizationRepository) GetOrganizationProvisioningCredential(
	ctx context.Context,
	organizationID int64,
	provider value.ProvisioningCredentialProvider,
) (*value.ProvisioningCredential, error) {
	credential, err := or.queries.GetOrganizationProvisioningCredential(ctx, sqlc.GetOrganizationProvisioningCredentialParams{
		OrganizationID: organizationID,
		Provider:       sqlc.Provider(provider),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &value.ProvisioningCredential{
		Provider: value.ProvisioningCredentialProvider(credential.Provider),
		Secret:   credential.Secret,
	}, nil
}

func (or *OrganizationRepository) AddOrganizationMember(ctx context.Context, organizationID int64, userID int64) error {
	err := or.queries.AddOrganizationMember(ctx, sqlc.AddOrganizationMemberParams{
		OrganizationID: organizationID,
		UserID:         userID,
	})

	return err
}

func (or *OrganizationRepository) RemoveOrganizationMember(ctx context.Context, organizationID int64, userID int64) error {
	err := or.queries.RemoveOrganizationMember(ctx, sqlc.RemoveOrganizationMemberParams{
		OrganizationID: organizationID,
		UserID:         userID,
	})

	return err
}

func (or *OrganizationRepository) CreateOrganizationInvite(ctx context.Context, organizationID int64, toEmail string, expiresAt time.Time) (*entity.OrganizationInvite, error) {
	invite, err := or.queries.CreateOrganizationInvite(ctx, sqlc.CreateOrganizationInviteParams{
		OrganizationID: organizationID,
		Email:          toEmail,
		ExpiresAt:      expiresAt,
	})

	if err != nil {
		return nil, err
	}

	return &entity.OrganizationInvite{
		Id:               invite.ID.String(),
		OrganizationId:   invite.OrganizationID,
		OrganizationName: invite.OrganizationName,
		Email:            invite.Email,
		ExpiresAt:        invite.ExpiresAt,
		CreatedAt:        invite.CreatedAt,
	}, nil
}

func (or *OrganizationRepository) GetOrganizationInviteById(ctx context.Context, inviteId string) (*entity.OrganizationInvite, error) {
	id, err := uuid.Parse(inviteId)
	if err != nil {
		return nil, err
	}
	invite, err := or.queries.GetOrganizationInviteById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.OrganizationInvite{
		Id:               invite.ID.String(),
		OrganizationId:   invite.OrganizationID,
		OrganizationName: invite.OrganizationName,
		Email:            invite.Email,
		ExpiresAt:        invite.ExpiresAt,
		CreatedAt:        invite.CreatedAt,
	}, nil
}

func (or *OrganizationRepository) GetOrganizationMembers(ctx context.Context, organizationID int64) ([]*entity.User, error) {
	rows, err := or.queries.GetOrganizationMembers(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return []*entity.User{}, nil
	}

	users := make([]*entity.User, 0, len(rows))
	for _, user := range rows {
		users = append(users, &entity.User{
			Id:           user.ID,
			BetterAuthId: user.BetterAuthID,
		})
	}
	return users, nil
}

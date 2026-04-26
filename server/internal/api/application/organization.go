package application

import (
	"context"
	"errors"
	"starliner.app/internal/api/conf"
	apiPort "starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	"time"
)

type OrganizationApplication struct {
	cfg                    *conf.Config
	crypto                 port.Crypto
	email                  apiPort.Email
	organizationRepository interfaces.OrganizationRepository
	teamRepository         interfaces.TeamRepository
	normalizationService   *coreService.NormalizerService
	organizationService    *service.OrganizationService
}

func NewOrganizationApplication(
	cfg *conf.Config,
	crypto port.Crypto,
	email apiPort.Email,
	organizationRepository interfaces.OrganizationRepository,
	teamRepository interfaces.TeamRepository,
	normalizationService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
) *OrganizationApplication {
	return &OrganizationApplication{
		cfg:                    cfg,
		crypto:                 crypto,
		email:                  email,
		organizationRepository: organizationRepository,
		teamRepository:         teamRepository,
		normalizationService:   normalizationService,
		organizationService:    organizationService,
	}
}

func (oa *OrganizationApplication) CreateOrganization(ctx context.Context, name string, ownerID int64) (*value.Organization, error) {
	organizationSlug, err := oa.normalizationService.FormatToDNS1123(name)
	if err != nil {
		return nil, err
	}

	org, err := oa.organizationRepository.CreateOrganization(ctx, name, organizationSlug, ownerID)
	if err != nil {
		return nil, err
	}

	err = oa.organizationRepository.AddOrganizationMember(ctx, org.Id, ownerID)
	if err != nil {
		return nil, err
	}

	team, err := oa.teamRepository.CreateTeam(ctx, org.Slug, org.Id)
	if err != nil {
		return nil, err
	}

	err = oa.teamRepository.AddTeamMember(ctx, team.Id, ownerID)
	if err != nil {
		return nil, err
	}

	return &value.Organization{
		Id:      org.Id,
		Name:    org.Name,
		Slug:    org.Slug,
		OwnerId: org.OwnerId,
	}, nil
}

func (oa *OrganizationApplication) GetUserOrganizations(ctx context.Context, userID int64) ([]*value.Organization, error) {
	organizations, err := oa.organizationRepository.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}
	return value.NewOrganizations(organizations), nil
}

func (oa *OrganizationApplication) GetProjectsForUser(ctx context.Context, userID int64, organizationID int64) ([]*value.Project, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	projects, err := oa.organizationRepository.GetUserProjects(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}
	return value.NewProjects(projects), nil
}

func (oa *OrganizationApplication) GetClustersForUser(ctx context.Context, userID int64, organizationID int64) ([]*value.Cluster, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	clusters, err := oa.organizationRepository.GetOrganizationClusters(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	return value.NewClusters(clusters), nil
}

func (oa *OrganizationApplication) UpsertHetznerCredential(ctx context.Context, userID int64, organizationID int64, apiKey string) error {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return err
	}

	apiKeyEncrypted, err := oa.crypto.Encrypt(apiKey)
	if err != nil {
		return err
	}

	err = oa.organizationRepository.UpsertProvisioningCredentials(ctx, organizationID, apiKeyEncrypted, value.HetznerCredential)
	return err
}

func (oa *OrganizationApplication) GetHetznerCredential(ctx context.Context, userID int64, organizationID int64) (*value.ProvisioningCredential, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	err = oa.organizationService.ValidateUserOrgOwner(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	c, err := oa.organizationRepository.GetOrganizationProvisioningCredential(ctx, organizationID, value.HetznerCredential)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	decrypted, err := oa.crypto.Decrypt(c.Secret)
	if err != nil {
		return nil, err
	}

	return &value.ProvisioningCredential{
		Provider: c.Provider,
		Secret:   decrypted,
	}, nil
}

func (oa *OrganizationApplication) GetInviteDetails(ctx context.Context, inviteID string) (*value.OrganizationInvite, error) {
	invite, err := oa.organizationRepository.GetOrganizationInviteById(ctx, inviteID)
	if err != nil {
		return nil, err
	}

	return value.NewOrganizationInvite(invite), nil
}

func (oa *OrganizationApplication) AcceptInvite(ctx context.Context, inviteID string, userID int64) error {
	invite, err := oa.organizationRepository.GetOrganizationInviteById(ctx, inviteID)
	if err != nil {
		return err
	}

	if time.Now().After(invite.ExpiresAt) {
		return errors.New("invite has expired")
	}

	err = oa.organizationRepository.AddOrganizationMember(ctx, invite.OrganizationId, userID)
	if err != nil {
		return err
	}

	org, err := oa.organizationRepository.GetOrganization(ctx, invite.OrganizationId)
	if err != nil {
		return err
	}

	// TODO: better to store default team explicitly, works as long as you can't change team slugs
	team, err := oa.teamRepository.GetTeamBySlug(ctx, org.Slug, org.Id)
	if err != nil {
		return err
	}

	return oa.teamRepository.AddTeamMember(ctx, team.Id, userID)
}

func (oa *OrganizationApplication) CreateAndSendEmailInvite(ctx context.Context, userID int64, organizationID int64, toEmail string, inviteUrlPrefix string) error {
	err := oa.organizationService.ValidateUserOrgOwner(ctx, organizationID, userID)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	invite, err := oa.organizationRepository.CreateOrganizationInvite(ctx, organizationID, expiresAt)
	if err != nil {
		return err
	}

	return oa.email.SendInvite(toEmail, apiPort.InviteData{
		OrganizationName: invite.OrganizationName,
		InviteLink:       inviteUrlPrefix + invite.Id,
	})
}

func (oa *OrganizationApplication) GetOrganizationMembers(ctx context.Context, userID int64, organizationID int64) ([]*value.User, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	members, err := oa.organizationRepository.GetOrganizationMembers(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return value.NewUsers(members), nil
}

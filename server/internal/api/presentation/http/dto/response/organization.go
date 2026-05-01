package response

import (
	"starliner.app/internal/api/domain/value"
)

type Organization struct {
	Id      int64  `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	OwnerId int64  `json:"owner_id" binding:"required"`
}

func NewOrganization(org *value.Organization) Organization {
	return Organization{
		Id:      org.Id,
		Name:    org.Name,
		Slug:    org.Slug,
		OwnerId: org.OwnerId,
	}
}

func NewOrganizations(orgs []*value.Organization) []Organization {
	res := make([]Organization, len(orgs))
	for i, org := range orgs {
		res[i] = NewOrganization(org)
	}
	return res
}

type OrganizationProvisioningCredential struct {
	Provider string `json:"provider" binding:"required, oneof=hetzner"`
	Secret   string `json:"secret" binding:"required"`
}

type GetOrganizationProvisioningCredentialResponse struct {
	Credential *OrganizationProvisioningCredential `json:"credential"`
}

type OrganizationInvite struct {
	Id               string `json:"id" binding:"required"`
	OrganizationId   int64  `json:"organization_id" binding:"required"`
	OrganizationName string `json:"organization_name" binding:"required"`
	Email            string `json:"email" binding:"required"`
	ExpiresAt        string `json:"expires_at" binding:"required"`
	CreatedAt        string `json:"created_at" binding:"required"`
}

func NewOrganizationInvite(invite *value.OrganizationInvite) OrganizationInvite {
	return OrganizationInvite{
		Id:               invite.Id,
		OrganizationId:   invite.OrganizationId,
		OrganizationName: invite.OrganizationName,
		Email:            invite.Email,
		ExpiresAt:        invite.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:        invite.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

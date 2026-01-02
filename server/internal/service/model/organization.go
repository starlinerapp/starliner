package model

import "starliner.app/internal/domain"

type Organization struct {
	Id      int64
	Name    string
	Slug    string
	OwnerId int64
}

func NewOrganization(o *domain.Organization) *Organization {
	return &Organization{
		Id:      o.Id,
		Name:    o.Name,
		Slug:    o.Slug,
		OwnerId: o.OwnerId,
	}
}

func NewOrganizations(os []*domain.Organization) []*Organization {
	organizations := make([]*Organization, len(os))
	for i, o := range os {
		organizations[i] = NewOrganization(o)
	}
	return organizations
}

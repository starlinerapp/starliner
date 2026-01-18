package value

import (
	"starliner.app/internal/api/domain/entity"
)

type Organization struct {
	Id      int64
	Name    string
	Slug    string
	OwnerId int64
}

func NewOrganization(o *entity.Organization) *Organization {
	return &Organization{
		Id:      o.Id,
		Name:    o.Name,
		Slug:    o.Slug,
		OwnerId: o.OwnerId,
	}
}

func NewOrganizations(os []*entity.Organization) []*Organization {
	organizations := make([]*Organization, len(os))
	for i, o := range os {
		organizations[i] = NewOrganization(o)
	}
	return organizations
}

package response

import "starliner.app/internal/service/model"

type Organization struct {
	Id      int64  `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	OwnerId int64  `json:"owner_id" binding:"required"`
}

func NewOrganization(org *model.Organization) Organization {
	return Organization{
		Id:      org.Id,
		Name:    org.Name,
		Slug:    org.Slug,
		OwnerId: org.OwnerId,
	}
}

func NewOrganizations(orgs []*model.Organization) []Organization {
	res := make([]Organization, len(orgs))
	for i, org := range orgs {
		res[i] = NewOrganization(org)
	}
	return res
}

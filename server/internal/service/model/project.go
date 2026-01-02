package model

import "starliner.app/internal/domain"

type Project struct {
	Id             int64
	Name           string
	Environments   []*Environment
	OrganizationId int64
}

func NewProject(p *domain.Project) *Project {
	environments := make([]*Environment, len(p.Environments))
	for i, e := range p.Environments {
		environments[i] = &Environment{
			Id:   e.Id,
			Slug: e.Slug,
			Name: e.Name,
		}
	}
	return &Project{
		Id:             p.Id,
		Name:           p.Name,
		Environments:   environments,
		OrganizationId: p.OrganizationId,
	}
}

func NewProjects(ps []*domain.Project) []*Project {
	projects := make([]*Project, len(ps))
	for i, p := range ps {
		projects[i] = NewProject(p)
	}
	return projects
}

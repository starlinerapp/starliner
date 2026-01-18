package value

import (
	"starliner.app/internal/api/domain/entity"
)

type Project struct {
	Id             int64
	Name           string
	Environments   []*Environment
	OrganizationId int64
	ClusterId      *int64
}

func NewProject(p *entity.Project) *Project {
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
		ClusterId:      p.ClusterId,
	}
}

func NewProjects(ps []*entity.Project) []*Project {
	projects := make([]*Project, len(ps))
	for i, p := range ps {
		projects[i] = NewProject(p)
	}
	return projects
}

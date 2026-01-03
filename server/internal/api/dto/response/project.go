package response

import (
	"starliner.app/internal/service/model"
)

type Project struct {
	Id           int64         `json:"id" binding:"required"`
	Name         string        `json:"name" binding:"required"`
	Environments []Environment `json:"environments" binding:"required"`
	ClusterId    *int64        `json:"clusterId" binding:"required"`
}

func NewProject(p *model.Project) Project {
	environments := make([]Environment, len(p.Environments))
	for i, e := range p.Environments {
		environments[i] = Environment{
			Id:   e.Id,
			Slug: e.Slug,
			Name: e.Name,
		}
	}
	return Project{
		Id:           p.Id,
		Name:         p.Name,
		Environments: environments,
		ClusterId:    p.ClusterId,
	}
}

func NewProjects(ps []*model.Project) []Project {
	projects := make([]Project, len(ps))
	for i, p := range ps {
		projects[i] = NewProject(p)
	}
	return projects
}

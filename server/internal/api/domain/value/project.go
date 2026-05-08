package value

import (
	"starliner.app/internal/api/domain/entity"
	"time"
)

type Project struct {
	Id           int64
	Name         string
	Environments []*Environment
	TeamId       int64
	TeamSlug     string
	ClusterId    *int64
	CreatedAt    time.Time
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
		Id:           p.Id,
		Name:         p.Name,
		Environments: environments,
		TeamId:       p.TeamId,
		TeamSlug:     p.TeamSlug,
		ClusterId:    p.ClusterId,
		CreatedAt:    p.CreatedAt,
	}
}

func NewProjects(ps []*entity.Project) []*Project {
	projects := make([]*Project, len(ps))
	for i, p := range ps {
		projects[i] = NewProject(p)
	}
	return projects
}

type ProjectCluster struct {
	ClusterId   int64
	ClusterName string
}

func NewProjectCluster(p *entity.ProjectCluster) *ProjectCluster {
	return &ProjectCluster{
		ClusterId:   p.ClusterId,
		ClusterName: p.ClusterName,
	}
}

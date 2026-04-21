package response

import (
	"time"

	"starliner.app/internal/api/domain/value"
)

type Project struct {
	Id           int64         `json:"id" binding:"required"`
	Name         string        `json:"name" binding:"required"`
	Environments []Environment `json:"environments" binding:"required"`
	ClusterId    *int64        `json:"clusterId" binding:"required"`
	TeamId       int64         `json:"teamId" binding:"required"`
	CreatedAt    time.Time     `json:"createdAt" binding:"required"`
}

func NewProject(p *value.Project) Project {
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
		TeamId:       p.TeamId,
		CreatedAt:    p.CreatedAt,
	}
}

func NewProjects(ps []*value.Project) []Project {
	projects := make([]Project, len(ps))
	for i, p := range ps {
		projects[i] = NewProject(p)
	}
	return projects
}

type ProjectCluster struct {
	ClusterId   int64  `json:"clusterId" binding:"required"`
	ClusterName string `json:"clusterName" binding:"required"`
}

func NewProjectCluster(p *value.ProjectCluster) ProjectCluster {
	return ProjectCluster{
		ClusterId:   p.ClusterId,
		ClusterName: p.ClusterName,
	}
}

type ProjectPreviewEnvironmentEnabled struct {
	Enabled bool `json:"enabled" binding:"required"`
}

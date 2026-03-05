package value

import (
	"starliner.app/internal/api/domain/entity"
)

type Deployments struct {
	Ingresses      []*IngressDeployment
	GitDeployments []*GitDeployment
	Databases      []*DatabaseDeployment
	Images         []*ImageDeployment
}

type PathType string

const (
	Prefix PathType = "Prefix"
	Exact  PathType = "Exact"
)

type GitDeployment struct {
	Id                    int64
	ServiceName           string
	Port                  string
	Status                string
	GitUrl                string
	ProjectRepositoryPath string
	DockerfilePath        string
	EnvVars               []*EnvVar
}

func NewGitDeployment(d *entity.GitDeployment) *GitDeployment {
	return &GitDeployment{
		Id:                    d.Id,
		ServiceName:           d.Name,
		Status:                d.Status,
		Port:                  d.Port,
		GitUrl:                d.GitUrl,
		ProjectRepositoryPath: d.ProjectRepositoryPath,
		DockerfilePath:        d.DockerfilePath,
		EnvVars:               mapEnvVars(d.EnvVars),
	}
}

func NewGitDeployments(ds []*entity.GitDeployment) []*GitDeployment {
	deployments := make([]*GitDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewGitDeployment(d)
	}
	return deployments
}

type IngressPath struct {
	Path        string
	PathType    PathType
	ServiceName string
}

type IngressHost struct {
	Host  string
	Paths []*IngressPath
}

type IngressDeployment struct {
	Id           int64
	ServiceName  string
	Status       string
	Port         string
	IngressHosts []*IngressHost
}

func NewIngressDeployment(d *entity.IngressDeployment) *IngressDeployment {
	return &IngressDeployment{
		Id:           d.Id,
		ServiceName:  d.Name,
		Status:       d.Status,
		Port:         d.Port,
		IngressHosts: mapIngressHosts(d.IngressHosts),
	}
}

func NewIngressDeployments(ds []*entity.IngressDeployment) []*IngressDeployment {
	deployments := make([]*IngressDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewIngressDeployment(d)
	}
	return deployments
}

func mapIngressHosts(hosts []*entity.IngressHost) []*IngressHost {
	if hosts == nil {
		return nil
	}

	result := make([]*IngressHost, 0, len(hosts))
	for _, h := range hosts {
		if h == nil {
			continue
		}

		result = append(result, &IngressHost{
			Host:  h.Host,
			Paths: mapIngressPaths(h.Paths),
		})
	}

	return result
}

func mapIngressPaths(paths []*entity.IngressPath) []*IngressPath {
	if paths == nil {
		return nil
	}

	result := make([]*IngressPath, 0, len(paths))
	for _, p := range paths {
		if p == nil {
			continue
		}

		result = append(result, &IngressPath{
			Path:        p.Path,
			PathType:    PathType(p.PathType),
			ServiceName: p.ServiceName,
		})
	}

	return result
}

type EnvVar struct {
	Name  string
	Value string
}

type ImageDeployment struct {
	Id          int64
	ServiceName string
	Status      string
	ImageName   string
	Tag         string
	Port        string
	EnvVars     []*EnvVar
}

func NewImageDeployment(d *entity.ImageDeployment) *ImageDeployment {
	return &ImageDeployment{
		Id:          d.Id,
		ServiceName: d.ServiceName,
		Status:      d.Status,
		ImageName:   d.ImageName,
		Tag:         d.Tag,
		Port:        d.Port,
		EnvVars:     mapEnvVars(d.EnvVars),
	}
}

func NewImageDeployments(ds []*entity.ImageDeployment) []*ImageDeployment {
	deployments := make([]*ImageDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewImageDeployment(d)
	}
	return deployments
}

func mapEnvVars(envVars []*entity.EnvVar) []*EnvVar {
	variables := make([]*EnvVar, len(envVars))
	for i, e := range envVars {
		variables[i] = &EnvVar{
			Name:  e.Name,
			Value: e.Value,
		}
	}
	return variables
}

type DatabaseDeployment struct {
	Id          int64
	ServiceName string
	Status      string
	Database    *string
	Username    *string
	Password    *string
	Port        string
}

func NewDatabaseDeployment(d *entity.DatabaseDeployment) *DatabaseDeployment {
	return &DatabaseDeployment{
		Id:          d.Id,
		ServiceName: d.ServiceName,
		Status:      d.Status,
		Database:    d.Database,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
	}
}

func NewDatabaseDeployments(ds []*entity.DatabaseDeployment) []*DatabaseDeployment {
	deployments := make([]*DatabaseDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewDatabaseDeployment(d)
	}
	return deployments
}

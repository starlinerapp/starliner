package value

import (
	"errors"

	"starliner.app/internal/api/domain/entity"
)

var ErrDeploymentNameAlreadyExists = errors.New("deployment name already exists")
var ErrIngressHostAlreadyExists = errors.New("ingress host already exists")

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
	InternalEndpoint      string
	Port                  string
	Status                string
	GitUrl                string
	ProjectRepositoryPath string
	DockerfilePath        string
	EnvVars               []*EnvVar
}

func NewGitDeployment(d *entity.GitDeployment, internalEndpoint string) *GitDeployment {
	return &GitDeployment{
		Id:                    d.Id,
		ServiceName:           d.Name,
		InternalEndpoint:      internalEndpoint,
		Status:                d.Status,
		Port:                  d.Port,
		GitUrl:                d.GitUrl,
		ProjectRepositoryPath: d.ProjectRepositoryPath,
		DockerfilePath:        d.DockerfilePath,
		EnvVars:               mapEnvVars(d.EnvVars),
	}
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
	Id               int64
	ServiceName      string
	InternalEndpoint string
	Status           string
	ImageName        string
	Tag              string
	Port             string
	VolumeSizeMB     *int32
	VolumeMountPath  *string
	EnvVars          []*EnvVar
}

func NewImageDeployment(d *entity.ImageDeployment, internalEndpoint string) *ImageDeployment {
	return &ImageDeployment{
		Id:               d.Id,
		ServiceName:      d.ServiceName,
		InternalEndpoint: internalEndpoint,
		Status:           d.Status,
		ImageName:        d.ImageName,
		Tag:              d.Tag,
		Port:             d.Port,
		VolumeSizeMB:     d.VolumeSizeMB,
		VolumeMountPath:  d.VolumeMountPath,
		EnvVars:          mapEnvVars(d.EnvVars),
	}
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
	Id               int64
	ServiceName      string
	InternalEndpoint string
	Status           string
	Database         *string
	Username         *string
	Password         *string
	Port             string
}

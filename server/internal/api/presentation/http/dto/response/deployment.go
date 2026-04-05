package response

import (
	"time"

	"starliner.app/internal/api/domain/value"
)

type Deployments struct {
	Ingresses     []IngressDeployment  `json:"ingresses" binding:"required"`
	Databases     []DatabaseDeployment `json:"databases" binding:"required"`
	Images        []ImageDeployment    `json:"images" binding:"required"`
	GitDeployment []GitDeployment      `json:"gitDeployments" binding:"required"`
}

type EnvVar struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type GitDeployment struct {
	Id                    int64    `json:"id" binding:"required"`
	ServiceName           string   `json:"serviceName" binding:"required"`
	InternalEndpoint      string   `json:"internalEndpoint" binding:"required"`
	Status                string   `json:"status" binding:"required"`
	Port                  string   `json:"port" binding:"required"`
	GitUrl                string   `json:"gitUrl" binding:"required"`
	ProjectRepositoryPath string   `json:"projectRepositoryPath" binding:"required"`
	DockerfilePath        string   `json:"dockerfilePath" binding:"required"`
	EnvVars               []EnvVar `json:"envVars" binding:"required"`
}

func NewGitDeployment(gitDeployment *value.GitDeployment) GitDeployment {
	return GitDeployment{
		Id:                    gitDeployment.Id,
		ServiceName:           gitDeployment.ServiceName,
		InternalEndpoint:      gitDeployment.InternalEndpoint,
		Status:                gitDeployment.Status,
		Port:                  gitDeployment.Port,
		GitUrl:                gitDeployment.GitUrl,
		ProjectRepositoryPath: gitDeployment.ProjectRepositoryPath,
		DockerfilePath:        gitDeployment.DockerfilePath,
		EnvVars:               mapEnvVarsFromValue(gitDeployment.EnvVars),
	}
}

func mapEnvVarsFromValue(envVars []*value.EnvVar) []EnvVar {
	variables := make([]EnvVar, len(envVars))
	for i, envVar := range envVars {
		variables[i] = EnvVar{
			Name:  envVar.Name,
			Value: envVar.Value,
		}
	}
	return variables
}

func NewGitDeployments(gitDeployments []*value.GitDeployment) []GitDeployment {
	result := make([]GitDeployment, 0, len(gitDeployments))
	for _, deployment := range gitDeployments {
		result = append(result, NewGitDeployment(deployment))
	}
	return result
}

type IngressPath struct {
	Path        string `json:"path" binding:"required"`
	PathType    string `json:"pathType" binding:"required,oneof=Prefix Exact"`
	ServiceName string `json:"serviceName" binding:"required"`
}

type IngressHost struct {
	Host  string        `json:"host" binding:"required"`
	Paths []IngressPath `json:"paths" binding:"required"`
}

type IngressDeployment struct {
	Id          int64         `json:"id" binding:"required"`
	ServiceName string        `json:"serviceName" binding:"required"`
	Status      string        `json:"status" binding:"required"`
	Port        string        `json:"port" binding:"required"`
	Hosts       []IngressHost `json:"hosts" binding:"required"`
}

func NewIngressDeployment(ingressDeployment *value.IngressDeployment) IngressDeployment {
	return IngressDeployment{
		Id:          ingressDeployment.Id,
		ServiceName: ingressDeployment.ServiceName,
		Status:      ingressDeployment.Status,
		Port:        ingressDeployment.Port,
		Hosts:       mapHostsFromValue(ingressDeployment.IngressHosts),
	}
}

func mapHostsFromValue(hosts []*value.IngressHost) []IngressHost {
	out := make([]IngressHost, 0, len(hosts))
	for _, h := range hosts {
		paths := make([]IngressPath, 0, len(h.Paths))
		for _, p := range h.Paths {
			paths = append(paths, IngressPath{
				Path:        p.Path,
				PathType:    string(p.PathType),
				ServiceName: p.ServiceName,
			})
		}
		out = append(out, IngressHost{
			Host:  h.Host,
			Paths: paths,
		})
	}
	return out
}

func NewIngressDeployments(ingressDeployments []*value.IngressDeployment) []IngressDeployment {
	result := make([]IngressDeployment, 0, len(ingressDeployments))
	for _, deployment := range ingressDeployments {
		result = append(result, NewIngressDeployment(deployment))
	}
	return result
}

type ImageDeployment struct {
	Id               int64    `json:"id" binding:"required"`
	ServiceName      string   `json:"serviceName" binding:"required"`
	InternalEndpoint string   `json:"internalEndpoint" binding:"required"`
	ImageName        string   `json:"imageName" binding:"required"`
	Tag              string   `json:"tag" binding:"required"`
	Status           string   `json:"status" binding:"required"`
	Port             string   `json:"port" binding:"required"`
	VolumeSizeMB     *int32   `json:"volumeSizeMB"`
	VolumeMountPath  *string  `json:"volumeMountPath"`
	EnvVars          []EnvVar `json:"envVars" binding:"required"`
}

func NewImageDeployment(imageDeployment *value.ImageDeployment) ImageDeployment {
	return ImageDeployment{
		Id:               imageDeployment.Id,
		ServiceName:      imageDeployment.ServiceName,
		InternalEndpoint: imageDeployment.InternalEndpoint,
		ImageName:        imageDeployment.ImageName,
		Tag:              imageDeployment.Tag,
		Status:           imageDeployment.Status,
		Port:             imageDeployment.Port,
		VolumeSizeMB:     imageDeployment.VolumeSizeMB,
		VolumeMountPath:  imageDeployment.VolumeMountPath,
		EnvVars:          mapEnvVarsFromValue(imageDeployment.EnvVars),
	}
}

func NewImageDeployments(imageDeployments []*value.ImageDeployment) []ImageDeployment {
	result := make([]ImageDeployment, 0, len(imageDeployments))
	for _, deployment := range imageDeployments {
		result = append(result, NewImageDeployment(deployment))
	}
	return result
}

type DatabaseDeployment struct {
	Id               int64   `json:"id" binding:"required"`
	ServiceName      string  `json:"serviceName" binding:"required"`
	InternalEndpoint string  `json:"internalEndpoint" binding:"required"`
	Status           string  `json:"status" binding:"required"`
	Database         *string `json:"database" binding:"required"`
	Username         *string `json:"username" binding:"required"`
	Password         *string `json:"password" binding:"required"`
	Port             string  `json:"port" binding:"required"`
}

func NewDatabaseDeployment(databaseDeployment *value.DatabaseDeployment) DatabaseDeployment {
	return DatabaseDeployment{
		Id:               databaseDeployment.Id,
		ServiceName:      databaseDeployment.ServiceName,
		InternalEndpoint: databaseDeployment.InternalEndpoint,
		Status:           databaseDeployment.Status,
		Database:         databaseDeployment.Database,
		Username:         databaseDeployment.Username,
		Password:         databaseDeployment.Password,
		Port:             databaseDeployment.Port,
	}
}

func NewDatabaseDeployments(databaseDeployments []*value.DatabaseDeployment) []DatabaseDeployment {
	result := make([]DatabaseDeployment, 0, len(databaseDeployments))
	for _, deployment := range databaseDeployments {
		result = append(result, NewDatabaseDeployment(deployment))
	}
	return result
}

func NewDeployments(deployments *value.Deployments) Deployments {
	return Deployments{
		Ingresses:     NewIngressDeployments(deployments.Ingresses),
		Databases:     NewDatabaseDeployments(deployments.Databases),
		Images:        NewImageDeployments(deployments.Images),
		GitDeployment: NewGitDeployments(deployments.GitDeployments),
	}
}

type GitDeploymentBuild struct {
	BuildId        int64     `json:"buildId" binding:"required"`
	DeploymentId   int64     `json:"deploymentId" binding:"required"`
	DeploymentName string    `json:"deploymentName" binding:"required"`
	Status         string    `json:"status" binding:"required"`
	GitUrl         string    `json:"gitUrl" binding:"required"`
	ProjectPath    string    `json:"projectPath" binding:"required"`
	DockerfilePath string    `json:"dockerfilePath" binding:"required"`
	CreatedAt      time.Time `json:"createdAt" binding:"required"`
}

func NewGitDeploymentBuild(build *value.GitDeploymentBuild) GitDeploymentBuild {
	return GitDeploymentBuild{
		BuildId:        build.BuildId,
		DeploymentId:   build.DeploymentId,
		DeploymentName: build.DeploymentName,
		Status:         string(build.Status),
		GitUrl:         build.GitUrl,
		ProjectPath:    build.ProjectPath,
		DockerfilePath: build.DockerfilePath,
		CreatedAt:      build.CreatedAt,
	}
}

func NewGitDeploymentBuilds(builds []*value.GitDeploymentBuild) []GitDeploymentBuild {
	result := make([]GitDeploymentBuild, 0, len(builds))
	for _, build := range builds {
		result = append(result, NewGitDeploymentBuild(build))
	}
	return result
}

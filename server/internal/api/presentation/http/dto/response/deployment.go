package response

import (
	"starliner.app/internal/api/domain/value"
)

type Deployments struct {
	Ingresses []IngressDeployment  `json:"ingresses" binding:"required"`
	Databases []DatabaseDeployment `json:"databases" binding:"required"`
	Images    []ImageDeployment    `json:"images" binding:"required"`
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

type EnvVar struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ImageDeployment struct {
	Id          int64    `json:"id" binding:"required"`
	ServiceName string   `json:"serviceName" binding:"required"`
	ImageName   string   `json:"imageName" binding:"required"`
	Tag         string   `json:"tag" binding:"required"`
	Status      string   `json:"status" binding:"required"`
	Port        string   `json:"port" binding:"required"`
	EnvVars     []EnvVar `json:"envVars" binding:"required"`
}

func NewImageDeployment(imageDeployment *value.ImageDeployment) ImageDeployment {
	return ImageDeployment{
		Id:          imageDeployment.Id,
		ServiceName: imageDeployment.ServiceName,
		ImageName:   imageDeployment.ImageName,
		Tag:         imageDeployment.Tag,
		Status:      imageDeployment.Status,
		Port:        imageDeployment.Port,
		EnvVars:     mapEnvVarsFromValue(imageDeployment.EnvVars),
	}
}

func NewImageDeployments(imageDeployments []*value.ImageDeployment) []ImageDeployment {
	result := make([]ImageDeployment, 0, len(imageDeployments))
	for _, deployment := range imageDeployments {
		result = append(result, NewImageDeployment(deployment))
	}
	return result
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

type DatabaseDeployment struct {
	Id          int64   `json:"id" binding:"required"`
	ServiceName string  `json:"serviceName" binding:"required"`
	Status      string  `json:"status" binding:"required"`
	Database    *string `json:"database" binding:"required"`
	Username    *string `json:"username" binding:"required"`
	Password    *string `json:"password" binding:"required"`
	Port        string  `json:"port" binding:"required"`
}

func NewDatabaseDeployment(databaseDeployment *value.DatabaseDeployment) DatabaseDeployment {
	return DatabaseDeployment{
		Id:          databaseDeployment.Id,
		ServiceName: databaseDeployment.ServiceName,
		Status:      databaseDeployment.Status,
		Database:    databaseDeployment.Database,
		Username:    databaseDeployment.Username,
		Password:    databaseDeployment.Password,
		Port:        databaseDeployment.Port,
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
		Ingresses: NewIngressDeployments(deployments.Ingresses),
		Databases: NewDatabaseDeployments(deployments.Databases),
		Images:    NewImageDeployments(deployments.Images),
	}
}

package response

import (
	"starliner.app/internal/api/domain/value"
)

type Deployments struct {
	Databases []DatabaseDeployment `json:"databases" binding:"required"`
	Images    []ImageDeployment    `json:"images" binding:"required"`
}

type ImageDeployment struct {
	Id          int64  `json:"id" binding:"required"`
	ServiceName string `json:"serviceName" binding:"required"`
	ImageName   string `json:"imageName" binding:"required"`
	Tag         string `json:"tag" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Port        string `json:"port" binding:"required"`
}

func NewImageDeployment(imageDeployment *value.ImageDeployment) ImageDeployment {
	return ImageDeployment{
		Id:          imageDeployment.Id,
		ServiceName: imageDeployment.ServiceName,
		ImageName:   imageDeployment.ImageName,
		Tag:         imageDeployment.Tag,
		Status:      imageDeployment.Status,
		Port:        imageDeployment.Port,
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
	Id       int64  `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Port     string `json:"port" binding:"required"`
}

func NewDatabaseDeployment(databaseDeployment *value.DatabaseDeployment) DatabaseDeployment {
	return DatabaseDeployment{
		Id:       databaseDeployment.Id,
		Name:     databaseDeployment.Name,
		Status:   databaseDeployment.Status,
		Username: databaseDeployment.Username,
		Password: databaseDeployment.Password,
		Port:     databaseDeployment.Port,
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
		Databases: NewDatabaseDeployments(deployments.Databases),
		Images:    NewImageDeployments(deployments.Images),
	}
}

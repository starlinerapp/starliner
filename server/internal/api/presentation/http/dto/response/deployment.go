package response

import (
	"starliner.app/internal/api/domain/value"
)

type Deployment struct {
	Id       int64  `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Port     string `json:"port" binding:"required"`
}

func NewDeployment(deployment *value.DatabaseDeployment) Deployment {
	return Deployment{
		Id:       deployment.Id,
		Name:     deployment.Name,
		Status:   deployment.Status,
		Username: deployment.Username,
		Password: deployment.Password,
		Port:     deployment.Port,
	}
}

func NewDeployments(deployments []*value.DatabaseDeployment) []Deployment {
	result := make([]Deployment, 0, len(deployments))
	for _, deployment := range deployments {
		result = append(result, NewDeployment(deployment))
	}
	return result
}

package response

import "starliner.app/internal/domain/value"

type Deployment struct {
	Name string `json:"name" binding:"required"`
}

func NewDeployment(deployment *value.Deployment) Deployment {
	return Deployment{Name: deployment.Name}
}

func NewDeployments(deployments []*value.Deployment) []Deployment {
	var result []Deployment
	for _, deployment := range deployments {
		result = append(result, NewDeployment(deployment))
	}
	return result
}

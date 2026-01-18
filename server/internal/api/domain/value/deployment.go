package value

import (
	"starliner.app/internal/api/domain/entity"
)

type Deployment struct {
	Id       int64
	Name     string
	Username string
	Password string
	Port     string
}

func NewDeployment(d *entity.Deployment) *Deployment {
	return &Deployment{
		Id:       d.Id,
		Name:     d.Name,
		Username: d.Username,
		Password: d.Password,
		Port:     d.Port,
	}
}

func NewDeployments(ds []*entity.Deployment) []*Deployment {
	deployments := make([]*Deployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewDeployment(d)
	}
	return deployments
}

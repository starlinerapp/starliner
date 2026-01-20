package value

import (
	"starliner.app/internal/api/domain/entity"
)

type DatabaseDeployment struct {
	Id       int64
	Name     string
	Username string
	Password string
	Port     string
}

func NewDatabaseDeployment(d *entity.DatabaseDeployment) *DatabaseDeployment {
	return &DatabaseDeployment{
		Id:       d.Id,
		Name:     d.Name,
		Username: d.Username,
		Password: d.Password,
		Port:     d.Port,
	}
}

func NewDatabaseDeployments(ds []*entity.DatabaseDeployment) []*DatabaseDeployment {
	deployments := make([]*DatabaseDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewDatabaseDeployment(d)
	}
	return deployments
}

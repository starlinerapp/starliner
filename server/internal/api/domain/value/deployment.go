package value

import (
	"starliner.app/internal/api/domain/entity"
)

type Deployments struct {
	Ingresses []*IngressDeployment
	Databases []*DatabaseDeployment
	Images    []*ImageDeployment
}

type IngressDeployment struct {
	Id          int64
	ServiceName string
	Status      string
	Port        string
}

func NewIngressDeployment(d *entity.IngressDeployment) *IngressDeployment {
	return &IngressDeployment{
		Id:          d.Id,
		ServiceName: d.Name,
		Status:      *d.Status,
		Port:        d.Port,
	}
}

func NewIngressDeployments(ds []*entity.IngressDeployment) []*IngressDeployment {
	deployments := make([]*IngressDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewIngressDeployment(d)
	}
	return deployments
}

type ImageDeployment struct {
	Id          int64
	ServiceName string
	Status      string
	ImageName   string
	Tag         string
	Port        string
}

func NewImageDeployment(d *entity.ImageDeployment) *ImageDeployment {
	return &ImageDeployment{
		Id:          d.Id,
		ServiceName: d.ServiceName,
		Status:      *d.Status,
		ImageName:   d.ImageName,
		Tag:         d.Tag,
		Port:        d.Port,
	}
}

func NewImageDeployments(ds []*entity.ImageDeployment) []*ImageDeployment {
	deployments := make([]*ImageDeployment, len(ds))
	for i, d := range ds {
		deployments[i] = NewImageDeployment(d)
	}
	return deployments
}

type DatabaseDeployment struct {
	Id       int64
	Name     string
	Status   string
	Username string
	Password string
	Port     string
}

func NewDatabaseDeployment(d *entity.DatabaseDeployment) *DatabaseDeployment {
	return &DatabaseDeployment{
		Id:       d.Id,
		Name:     d.Name,
		Status:   *d.Status,
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

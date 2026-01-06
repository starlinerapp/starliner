package value

import (
	"starliner.app/internal/domain/entity"
)

type Cluster struct {
	Id             int64
	Name           string
	Status         entity.ClusterStatus
	IPv4Address    *string
	OrganizationId int64
}

func NewCluster(c *entity.Cluster) *Cluster {
	return &Cluster{
		Id:             c.Id,
		Name:           c.Name,
		Status:         c.Status,
		IPv4Address:    c.IPv4Address,
		OrganizationId: c.OrganizationId,
	}
}

func NewClusters(cs []*entity.Cluster) []*Cluster {
	clusters := make([]*Cluster, len(cs))
	for i, c := range cs {
		clusters[i] = NewCluster(c)
	}
	return clusters
}

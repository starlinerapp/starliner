package model

import "starliner.app/internal/domain"

type Cluster struct {
	Id             int64
	Name           string
	Status         domain.ClusterStatus
	IPv4Address    *string
	OrganizationId int64
}

func NewCluster(c *domain.Cluster) *Cluster {
	return &Cluster{
		Id:             c.Id,
		Name:           c.Name,
		Status:         c.Status,
		IPv4Address:    c.IPv4Address,
		OrganizationId: c.OrganizationId,
	}
}

func NewClusters(cs []*domain.Cluster) []*Cluster {
	clusters := make([]*Cluster, len(cs))
	for i, c := range cs {
		clusters[i] = NewCluster(c)
	}
	return clusters
}

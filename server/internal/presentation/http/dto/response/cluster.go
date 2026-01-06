package response

import (
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/value"
)

type Cluster struct {
	Id             int64                `json:"id" binding:"required"`
	Name           string               `json:"name" binding:"required"`
	Status         entity.ClusterStatus `json:"status" binding:"required,oneof=pending running deleted"`
	IPv4Address    *string              `json:"ipv4Address" binding:"required"`
	OrganizationId int64                `json:"organizationId" binding:"required"`
}

func NewCluster(cluster *value.Cluster) Cluster {
	return Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         cluster.Status,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
	}
}

func NewClusters(clusters []*value.Cluster) []Cluster {
	res := make([]Cluster, len(clusters))
	for i, c := range clusters {
		res[i] = NewCluster(c)
	}
	return res
}

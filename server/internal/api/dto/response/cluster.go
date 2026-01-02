package response

import (
	"starliner.app/internal/domain"
	"starliner.app/internal/service/model"
)

type Cluster struct {
	Id             int64                `json:"id" binding:"required"`
	Name           string               `json:"name" binding:"required"`
	Status         domain.ClusterStatus `json:"status" binding:"required,oneof=pending running deleted"`
	IPv4Address    *string              `json:"ipv4Address" binding:"required"`
	OrganizationId int64                `json:"organizationId" binding:"required"`
}

func NewCluster(cluster *model.Cluster) Cluster {
	return Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         cluster.Status,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
	}
}

func NewClusters(clusters []*model.Cluster) []Cluster {
	res := make([]Cluster, len(clusters))
	for i, c := range clusters {
		res[i] = NewCluster(c)
	}
	return res
}

package response

import (
	"time"

	"starliner.app/internal/api/domain/value"
)

type ClusterStatus string

const (
	ClusterStatusPending ClusterStatus = "pending"
	ClusterStatusRunning ClusterStatus = "running"
	ClusterStatusDeleted ClusterStatus = "deleted"
)

type Cluster struct {
	Id             int64         `json:"id" binding:"required"`
	Name           string        `json:"name" binding:"required"`
	Status         ClusterStatus `json:"status" binding:"required,oneof=pending running deleted"`
	TeamSlug       *string       `json:"teamSlug" binding:"required"`
	User           string        `json:"user" binding:"required"`
	IPv4Address    *string       `json:"ipv4Address" binding:"required"`
	OrganizationId int64         `json:"organizationId" binding:"required"`
	CreatedAt      time.Time     `json:"createdAt" binding:"required"`
	ServerType     string        `json:"serverType" binding:"required"`
}

func NewClusters(clusters []*value.Cluster) []Cluster {
	res := make([]Cluster, len(clusters))
	for i, c := range clusters {
		res[i] = NewCluster(c)
	}
	return res
}

func NewCluster(cluster *value.Cluster) Cluster {
	return Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         ClusterStatus(cluster.Status),
		TeamSlug:       cluster.TeamSlug,
		User:           cluster.User,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
		CreatedAt:      cluster.CreatedAt,
		ServerType:     string(cluster.ServerType),
	}
}

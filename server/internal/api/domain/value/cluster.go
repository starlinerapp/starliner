package value

import (
	"starliner.app/internal/api/domain/entity"
	"time"
)

type ClusterStatus string

const (
	ClusterStatusPending ClusterStatus = "pending"
	ClusterStatusRunning ClusterStatus = "running"
	ClusterStatusDeleted ClusterStatus = "deleted"
)

type ServerType string

const (
	ServerTypeCX23  ServerType = "cx23"
	ServerTypeCCX33 ServerType = "ccx33"
)

type Cluster struct {
	Id             int64
	Name           string
	Status         ClusterStatus
	IPv4Address    *string
	OrganizationId int64
	CreatedAt      time.Time
	ServerType     ServerType
}

func NewClusters(cs []*entity.Cluster) []*Cluster {
	clusters := make([]*Cluster, len(cs))
	for i, c := range cs {
		clusters[i] = NewCluster(c)
	}
	return clusters
}

func NewCluster(c *entity.Cluster) *Cluster {
	return &Cluster{
		Id:             c.Id,
		Name:           c.Name,
		Status:         mapStatus(c.Status),
		IPv4Address:    c.IPv4Address,
		OrganizationId: c.OrganizationId,
		CreatedAt:      c.CreatedAt,
		ServerType:     ServerType(c.ServerType),
	}
}

func mapStatus(s entity.ClusterStatus) ClusterStatus {
	switch s {
	case entity.ClusterPending:
		return ClusterStatusPending
	case entity.ClusterRunning:
		return ClusterStatusRunning
	case entity.ClusterDeleted:
		return ClusterStatusDeleted
	default:
		return "unknown"
	}
}

package model

import "starliner.app/internal/domain"

type Cluster struct {
	Id             int64
	Name           string
	Status         domain.ClusterStatus
	IPv4Address    *string
	OrganizationId int64
}

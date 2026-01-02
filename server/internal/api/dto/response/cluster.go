package response

import "starliner.app/internal/domain"

type Cluster struct {
	Id             int64                `json:"id" binding:"required"`
	Name           string               `json:"name" binding:"required"`
	Status         domain.ClusterStatus `json:"status" binding:"required,oneof=pending running deleted"`
	IPv4Address    *string              `json:"ipv4Address" binding:"required"`
	OrganizationId int64                `json:"organizationId" binding:"required"`
}

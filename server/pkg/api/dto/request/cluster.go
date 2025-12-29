package request

type CreateCluster struct {
	Name           string `json:"name" binding:"required"`
	OrganizationID int64  `json:"organizationId" binding:"required"`
}

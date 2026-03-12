package request

type CreateCluster struct {
	Name           string `json:"name" binding:"required"`
	ServerType     string `json:"serverType" binding:"required,oneof=cx23 ccx33"`
	OrganizationID int64  `json:"organizationId" binding:"required"`
}

package request

type CreateEnvironment struct {
	Name                string `json:"name" binding:"required"`
	ProjectID           int64  `json:"projectId" binding:"required"`
	OrganizationID      int64  `json:"organizationId" binding:"required"`
	SourceEnvironmentID *int64 `json:"sourceEnvironmentId"`
}

type UpdateEnvironmentConnectBranch struct {
	Branch string `json:"branch" binding:"required"`
}

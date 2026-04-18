package request

type CreateEnvironment struct {
	Name                string `json:"name" binding:"required"`
	ProjectID           int64  `json:"project_id" binding:"required"`
	OrganizationID      int64  `json:"organization_id" binding:"required"`
	SourceEnvironmentID *int64 `json:"source_environment_id"`
}

type UpdateEnvironmentConnectBranch struct {
	Branch string `json:"branch" binding:"required"`
}

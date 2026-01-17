package request

type CreateProject struct {
	Name           string `json:"name" binding:"required"`
	OrganizationId int64  `json:"organization_id" binding:"required"`
	ClusterId      int64  `json:"cluster_id" binding:"required"`
}

type UpdateProjectName struct {
	Name string `json:"name" binding:"required"`
}

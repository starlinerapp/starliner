package request

type CreateGithubApp struct {
	OrganizationId int64 `json:"organizationId" binding:"required"`
	InstallationId int64 `json:"installationId" binding:"required"`
}

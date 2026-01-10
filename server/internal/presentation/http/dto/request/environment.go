package request

type CreateEnvironment struct {
	Name           string `json:"name" binding:"required"`
	ProjectID      int64  `json:"project_id" binding:"required"`
	OrganizationID int64  `json:"organization_id" binding:"required"`
}

type Database string

const (
	Postgres Database = "postgres"
)

type DeployDatabase struct {
	Database Database `json:"database" binding:"required,oneof=postgres"`
}

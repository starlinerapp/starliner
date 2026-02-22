package request

type Database string

const (
	Postgres Database = "postgres"
)

type DeployApplication struct {
	EnvironmentId int64 `json:"environmentId" binding:"required"`
}

type DeployDatabase struct {
	EnvironmentId int64    `json:"environmentId" binding:"required"`
	Database      Database `json:"database" binding:"required,oneof=postgres"`
}

type DeployIngress struct {
	EnvironmentId int64 `json:"environmentId" binding:"required"`
}

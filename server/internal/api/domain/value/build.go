package value

import "time"

type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusBuilding BuildStatus = "building"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failure"
)

type GitDeploymentBuild struct {
	BuildId        int64
	DeploymentId   int64
	DeploymentName string
	CommitHash     *string
	Source         string
	Status         BuildStatus
	GitUrl         string
	ProjectPath    string
	DockerfilePath string
	CreatedAt      time.Time
}

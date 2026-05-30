package value

import "time"

type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusBuilding BuildStatus = "building"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failure"
)

type Arg struct {
	Name  string
	Value string
}

type GitDeploymentBuild struct {
	BuildId                 int64
	DeploymentId            int64
	DeploymentName          string
	DeploymentRolloutStatus string
	CommitHash              *string
	Source                  string
	Status                  BuildStatus
	CreatedAt               time.Time
}

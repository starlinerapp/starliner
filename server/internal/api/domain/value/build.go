package value

type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusBuilding BuildStatus = "building"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failed"
)

type Build struct {
	Id           int64
	DeploymentId int64
	Status       BuildStatus
	Logs         string
}

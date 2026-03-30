package value

type TriggerBuild struct {
	BuildId        int64
	DeploymentId   int64
	ImageName      string
	GitUrl         string
	AccessToken    string
	RootDirectory  string
	DockerfilePath string
}

type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusBuilding BuildStatus = "building"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failure"
)

type BuildCompleted struct {
	BuildId          int64
	DeploymentId     int64
	ImageRegistryUrl string
	BuildStatus      BuildStatus
	ImageName        string
	Tag              string
	Logs             string
}

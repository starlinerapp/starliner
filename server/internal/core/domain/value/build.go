package value

type Arg struct {
	Name  string
	Value string
}

type TriggerBuild struct {
	BuildId        int64
	DeploymentId   int64
	ImageName      string
	GitUrl         string
	AccessToken    string
	RootDirectory  string
	DockerfilePath string
	Args           []*Arg
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
	CommitHash       *string
	Tag              *string
	Logs             string
}

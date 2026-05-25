package value

type Arg struct {
	Name  string
	Value string
}

type TriggerBuild struct {
	BuildId        int64
	DeploymentId   int64
	CorrelationId  *string
	ImageName      string
	GitUrl         string
	BranchName     string
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
	CorrelationId    *string
	ImageRegistryUrl string
	BuildStatus      BuildStatus
	ImageName        *string
	CommitHash       *string
	Tag              *string
	Logs             string
}

type BuildLogChunk struct {
	BuildId int64
	Data    []byte
	End     bool
}


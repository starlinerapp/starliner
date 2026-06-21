package value

type Arg struct {
	Name  string
	Value string
}

type TriggerBuild struct {
	BuildId           int64
	DeploymentId      int64
	OrganizationId    int64
	ImageName         string
	GitUrl            string
	BranchName        string
	AccessToken       string
	RegistryPushToken string
	RootDirectory     string
	DockerfilePath    string
	Args              []*Arg
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

type RunnerBuildJob struct {
	RunnerId          int64
	BuildId           int64
	DeploymentId      int64
	ImageName         string
	ImageRegistryUrl  string
	GitUrl            string
	BranchName        string
	AccessToken       string
	RegistryPushToken string
	RootDirectory     string
	DockerfilePath    string
	Args              []*Arg
}

type RunnerBuildResult struct {
	BuildId      int64
	DeploymentId int64
	CommitHash   string
	Tag          string
	ImageName    string
	Logs         string
	Status       BuildStatus
}

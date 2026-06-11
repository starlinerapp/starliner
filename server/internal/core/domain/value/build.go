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

type BuildFailureStage string

const (
	BuildFailureStageClone BuildFailureStage = "clone"
	BuildFailureStageBuild BuildFailureStage = "build"
)

type BuildSucceeded struct {
	CorrelationId    string `json:"correlationId"`
	BuildId          int64  `json:"buildId"`
	DeploymentId     int64  `json:"deploymentId"`
	ImageRegistryUrl string `json:"imageRegistryUrl"`
	ImageName        string `json:"imageName"`
	CommitHash       string `json:"commitHash"`
	Tag              string `json:"tag"`
	Logs             string `json:"logs"`
}
type BuildFailed struct {
	CorrelationId string            `json:"correlationId"`
	BuildId       int64             `json:"buildId"`
	DeploymentId  int64             `json:"deploymentId"`
	ImageName     string            `json:"imageName"`
	GitUrl        string            `json:"gitUrl"`
	Stage         BuildFailureStage `json:"stage"`
	Logs          string            `json:"logs"`
}

type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusBuilding BuildStatus = "building"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failure"
)

type BuildLogChunk struct {
	BuildId int64
	Data    []byte
	End     bool
}

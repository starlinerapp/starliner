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
	ImageRegistryUrl string
	BuildStatus      BuildStatus
	ImageName        *string
	CommitHash       *string
	Tag              *string
	Logs             string
}

// BuildLogChunk is a single chunk of streaming build-log output. End is true on
// the final message emitted for a build, signaling that no further chunks will
// arrive on the stream.
type BuildLogChunk struct {
	BuildId int64
	Data    []byte
	End     bool
}

package value

type TriggerBuild struct {
	DeploymentId   int64
	ImageName      string
	GitUrl         string
	RootDirectory  string
	DockerfilePath string
}

type BuildCompleted struct {
	DeploymentId     int64
	ImageRegistryUrl string
	ImageName        string
	Tag              string
}

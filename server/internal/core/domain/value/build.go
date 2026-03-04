package value

type TriggerBuild struct {
	Id             int64
	ImageName      string
	GitUrl         string
	RootDirectory  string
	DockerfilePath string
}

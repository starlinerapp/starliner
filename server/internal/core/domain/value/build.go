package value

type TriggerBuild struct {
	Id             string
	Organization   string
	Project        string
	Service        string
	S3Key          string
	RootDirectory  string
	DockerfilePath string
}

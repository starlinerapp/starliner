package value

type BuildMessage struct {
	Id             string
	Organization   string
	Project        string
	Service        string
	S3Key          string
	RootDirectory  string
	DockerfilePath string
}

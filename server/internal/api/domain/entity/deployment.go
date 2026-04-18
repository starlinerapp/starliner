package entity

type PathType string

const (
	Prefix PathType = "Prefix"
	Exact  PathType = "Exact"
)

type IngressPath struct {
	Path        string
	PathType    PathType
	ServiceName string
}

type IngressHost struct {
	Host  string
	Paths []*IngressPath
}

type IngressDeployment struct {
	Id            int64
	EnvironmentId int64
	Status        string
	Name          string
	Port          string
	IngressHosts  []*IngressHost
}

type EnvVar struct {
	Name  string
	Value string
}

type Arg struct {
	Name  string
	Value string
}

type ImageDeployment struct {
	Id              int64
	Status          string
	ServiceName     string
	ImageName       string
	Tag             string
	Port            string
	EnvironmentId   int64
	VolumeSizeMiB   *int32
	VolumeMountPath *string
	EnvVars         []*EnvVar
}

type DatabaseDeployment struct {
	Id            int64
	ServiceName   string
	Status        string
	Database      *string
	Username      *string
	Password      *string
	Port          string
	EnvironmentId int64
}

type GitDeployment struct {
	Id                    int64
	Name                  string
	Status                string
	Port                  string
	EnvironmentId         int64
	GitUrl                string
	ProjectRepositoryPath string
	DockerfilePath        string
	EnvVars               []*EnvVar
}

type Deployment struct {
	Id            int64
	Name          string
	Port          string
	Namespace     string
	EnvironmentId int64
}

type DeploymentWithKubeconfig struct {
	Deployment Deployment
	Kubeconfig *string
}

type IngressHostDeployment struct {
	Host         string
	DeploymentId int64
}

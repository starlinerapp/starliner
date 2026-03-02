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
	Status        *string
	Name          string
	Port          string
	IngressHosts  []*IngressHost
}

type EnvVar struct {
	Name  string
	Value string
}

type ImageDeployment struct {
	Id            int64
	Status        *string
	ServiceName   string
	ImageName     string
	Tag           string
	Port          string
	EnvironmentId int64
	EnvVars       []*EnvVar
}

type DatabaseDeployment struct {
	Id            int64
	Name          string
	Status        *string
	Username      string
	Password      string
	Port          string
	EnvironmentId int64
}

type Deployment struct {
	Id            int64
	Name          string
	Port          string
	EnvironmentId int64
}

type DeploymentWithKubeconfig struct {
	Deployment Deployment
	Kubeconfig *string
}

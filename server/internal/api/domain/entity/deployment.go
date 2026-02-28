package entity

type IngressDeployment struct {
	Id            int64
	Status        *string
	Name          string
	Port          string
	EnvironmentId int64
}

type ImageDeployment struct {
	Id            int64
	Status        *string
	ServiceName   string
	ImageName     string
	Tag           string
	Port          string
	EnvironmentId int64
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

package entity

type Deployment struct {
	Id            int64
	Name          string
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

type DeploymentWithKubeconfig struct {
	Deployment Deployment
	Kubeconfig *string
}

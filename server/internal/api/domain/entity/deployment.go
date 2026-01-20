package entity

type Deployment struct {
	Name          string
	Port          string
	EnvironmentId int64
}

type DatabaseDeployment struct {
	Id            int64
	Name          string
	Username      string
	Password      string
	Port          string
	EnvironmentId int64
}

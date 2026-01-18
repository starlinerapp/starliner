package entity

type Deployment struct {
	Id            int64
	Name          string
	Username      string
	Password      string
	Port          string
	EnvironmentId int64
}

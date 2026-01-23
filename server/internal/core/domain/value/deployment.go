package value

type Deployment struct {
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
}

type DeploymentDeleted struct {
	DeploymentId int64
}

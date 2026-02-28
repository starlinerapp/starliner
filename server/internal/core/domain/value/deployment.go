package value

type ImageDeployment struct {
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
	ImageRepository  string
	ImageTag         string
	Port             int
}

type Deployment struct {
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
}

type IngressDeployment struct {
	HostName         string
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
}

type DeploymentDeleted struct {
	DeploymentId int64
}

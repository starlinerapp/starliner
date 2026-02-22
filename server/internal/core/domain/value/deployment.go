package value

type ApplicationDeployment struct {
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
	ImageRepository  string
	ImageTag         string
	Port             int
}

type DatabaseDeployment struct {
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
}

type IngressDeployment struct {
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
}

type DeploymentDeleted struct {
	DeploymentId int64
}

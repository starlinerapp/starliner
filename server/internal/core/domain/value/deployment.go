package value

type EnvVar struct {
	Name  string
	Value string
}

type ImageDeployment struct {
	DeploymentId     int64
	DeploymentName   string
	Namespace        string
	KubeconfigBase64 string
	ImageName        string
	ImageTag         string
	Port             int
	VolumeSizeMiB    *int32
	VolumeMountPath  *string
	EnvVars          []*EnvVar
}

type DatabaseDeployment struct {
	DeploymentId int64
	DbName       string
	Username     string
	Password     string
}

type Deployment struct {
	Namespace        string
	DeploymentId     int64
	DeploymentName   string
	KubeconfigBase64 string
}

type PathType string

const (
	Prefix PathType = "Prefix"
	Exact  PathType = "Exact"
)

type IngressPath struct {
	Path        string
	PathType    PathType
	ServiceName string
	ServicePort int
}

type IngressHost struct {
	Host  string
	Paths []IngressPath
}

type IngressDeployment struct {
	DeploymentId     int64
	DeploymentName   string
	Namespace        string
	KubeconfigBase64 string
	IngressHosts     []IngressHost
}

type DeploymentDeleted struct {
	DeploymentId int64
}

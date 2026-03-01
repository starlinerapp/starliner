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
	KubeconfigBase64 string
	IngressHosts     []IngressHost
}

type DeploymentDeleted struct {
	DeploymentId int64
}

package port

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

type DeployIngressArgs struct {
	Namespace        string
	ReleaseName      string
	KubeconfigBase64 string
	Hosts            []IngressHost
}

type EnvVar struct {
	Name  string
	Value string
}

type DeployImageArgs struct {
	Namespace        string
	ReleaseName      string
	KubeconfigBase64 string
	ImageRepository  string
	ImageTag         string
	Port             int
	EnvVars          []*EnvVar
}

type Deploy interface {
	DeployImage(args *DeployImageArgs) error

	DeployCloudNativePg(namespace string, releaseName string, kubeconfigBase64 string) error

	DeployPostgres(namespace string, releaseName string, kubeconfigBase64 string) error
	DeleteDeployment(namespace string, releaseName string, kubeconfigBase64 string) error

	DeployIngress(args *DeployIngressArgs) error
}

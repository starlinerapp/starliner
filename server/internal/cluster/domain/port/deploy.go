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
	ReleaseName      string
	KubeconfigBase64 string
	Hosts            []IngressHost
}

type EnvVar struct {
	Name  string
	Value string
}

type DeployImageArgs struct {
	ReleaseName      string
	KubeconfigBase64 string
	ImageRepository  string
	ImageTag         string
	Port             int
	EnvVars          []*EnvVar
}

type Deploy interface {
	DeployImage(args *DeployImageArgs) error

	DeployCloudNativePg(releaseName string, kubeconfigBase64 string) error

	DeployPostgres(releaseName string, kubeconfigBase64 string) error
	DeleteDeployment(releaseName string, kubeconfigBase64 string) error

	DeployIngress(args *DeployIngressArgs) error
}

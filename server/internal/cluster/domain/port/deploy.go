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

type Deploy interface {
	DeployImage(
		releaseName string,
		kubeconfigBase64 string,
		imageRepository string,
		imageTag string,
		port int,
	) error

	DeployCloudNativePg(releaseName string, kubeconfigBase64 string) error

	DeployPostgres(releaseName string, kubeconfigBase64 string) error
	DeletePostgres(releaseName string, kubeconfigBase64 string) error

	DeployIngress(args *DeployIngressArgs) error
}

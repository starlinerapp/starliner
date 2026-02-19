package port

type IngressPath struct {
	Path        string
	PathType    string
	ServiceName string
	ServicePort int
}

type IngressHost struct {
	Host  string
	Paths []IngressPath
}

type DeployIngressArgs struct {
	ReleaseName    string
	KubeconfigPath string
	Hosts          []IngressHost
}

type Deploy interface {
	DeployCloudNativePg(releaseName string, kubeconfigPath string) error

	DeployPostgres(releaseName string, kubeconfigPath string) error
	DeletePostgres(releaseName string, kubeconfigPath string) error

	DeployIngress(args *DeployIngressArgs) error
}

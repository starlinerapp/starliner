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
	ReleaseName      string
	KubeconfigBase64 string
	Hosts            []IngressHost
}

type Deploy interface {
	DeployCloudNativePg(releaseName string, kubeconfigBase64 string) error

	DeployPostgres(releaseName string, kubeconfigBase64 string) error
	DeletePostgres(releaseName string, kubeconfigBase64 string) error

	DeployIngress(args *DeployIngressArgs) error
}

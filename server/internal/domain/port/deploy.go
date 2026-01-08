package port

type Deploy interface {
	DeployNginx(ip string, kubeconfigPath string) error
}

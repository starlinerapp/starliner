package port

type Deploy interface {
	DeployPostgres(releaseName string, kubeconfigPath string) error
}

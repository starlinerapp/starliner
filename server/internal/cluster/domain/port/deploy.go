package port

type Deploy interface {
	DeployPostgres(releaseName string, kubeconfigPath string) error
	DeletePostgres(releaseName string, kubeconfigPath string) error
}

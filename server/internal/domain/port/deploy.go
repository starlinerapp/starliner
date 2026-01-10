package port

type Deploy interface {
	DeployPostgres(kubeconfigPath string) error
}

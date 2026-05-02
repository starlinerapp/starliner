package port

type Install interface {
	InstallK3s(clusterId int64, ip string, privateKey []byte) (kubeconfig string, logs string, err error)
}

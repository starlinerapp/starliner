package port

type Install interface {
	InstallK3s(ip string, privateKey []byte) (kubeconfig string, err error)
}

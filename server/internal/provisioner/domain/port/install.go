package port

type Install interface {
	InstallK3s(provisioningId string, ip string, privateKey []byte) (kubeconfig string, err error)
}

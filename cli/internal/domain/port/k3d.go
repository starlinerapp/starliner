package port

type K3dClient interface {
	Install() error
	Start() error
	Stop() error
}

package port

type WgClient interface {
	Install() error
}

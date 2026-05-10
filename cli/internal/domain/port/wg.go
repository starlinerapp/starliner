package port

type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

type WgClient interface {
	Install() error
	GenerateKeyPair() (keyPair *KeyPair, err error)
}

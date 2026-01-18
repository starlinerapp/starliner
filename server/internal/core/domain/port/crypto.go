package port

type Crypto interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
	GenerateKeyPair() (publicKey []byte, privateKey []byte, err error)
	EncodePrivateKeyToPEM(privateKey []byte) ([]byte, error)
}

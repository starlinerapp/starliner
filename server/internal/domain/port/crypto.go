package port

type Crypto interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
	EncodePrivateKeyToPEM(privateKey []byte) ([]byte, error)
}

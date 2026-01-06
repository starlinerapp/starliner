package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"starliner.app/internal/conf"
)

type Crypto struct {
	cfg *conf.Config
}

func NewCrypto(cfg *conf.Config) *Crypto {
	return &Crypto{cfg: cfg}
}

func (c *Crypto) Encrypt(plaintext string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(c.cfg.EncryptionKeyBase64)
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}

	encrypted, err := encryptAES(plaintext, encryptionKey)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func (c *Crypto) Decrypt(ciphertext string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(c.cfg.EncryptionKeyBase64)
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}

	decrypted, err := decryptAES(ciphertext, encryptionKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func (c *Crypto) EncodePrivateKeyToPEM(privateKey ed25519.PrivateKey) ([]byte, error) {
	block, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return nil, err
	}

	pemBytes := pem.EncodeToMemory(block)
	if pemBytes == nil {
		return nil, err
	}
	return pemBytes, nil
}

func encryptAES(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(ciphertextB64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, data := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

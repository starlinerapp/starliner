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
	"starliner.app/internal/core/conf"
	"starliner.app/internal/core/domain/port"
)

type Crypto struct {
	cfg conf.CryptoConfig
}

func NewCrypto(cfg conf.CryptoConfig) port.Crypto {
	return &Crypto{cfg: cfg}
}

func (c *Crypto) Encrypt(plaintext string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(c.cfg.GetEncryptionKeyBase64())
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
	encryptionKey, err := base64.StdEncoding.DecodeString(c.cfg.GetEncryptionKeyBase64())
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}

	decrypted, err := decryptAES(ciphertext, encryptionKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func (c *Crypto) GenerateKeyPair() (publicKey []byte, privateKey []byte, err error) {
	return ed25519.GenerateKey(rand.Reader)
}

func (c *Crypto) EncodePrivateKeyToPEM(privateKey []byte) ([]byte, error) {
	block, err := ssh.MarshalPrivateKey(ed25519.PrivateKey(privateKey), "")
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

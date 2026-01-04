package service

import (
	"encoding/base64"
	"fmt"
	"starliner.app/internal/config"
	"starliner.app/internal/crypto"
)

type CryptoService struct {
	cfg *config.Config
}

func NewCryptoService(cfg *config.Config) *CryptoService {
	return &CryptoService{cfg: cfg}
}

func (cs *CryptoService) Encrypt(plaintext string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(cs.cfg.EncryptionKeyBase64)
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}

	encrypted, err := crypto.Encrypt(plaintext, encryptionKey)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func (cs *CryptoService) Decrypt(ciphertext string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(cs.cfg.EncryptionKeyBase64)
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}

	decrypted, err := crypto.Decrypt(ciphertext, encryptionKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

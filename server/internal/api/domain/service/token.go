package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"
)

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (ts *TokenService) GenerateToken() (string, error) {
	raw := make([]byte, 20)

	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw)

	return strings.ToUpper(token), nil
}

func (ts *TokenService) HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

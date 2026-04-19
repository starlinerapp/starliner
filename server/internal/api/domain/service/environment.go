package service

import (
	"context"
	"errors"
	"math/rand"
	"starliner.app/internal/api/domain/repository/interface"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

type EnvironmentService struct {
	environmentRepository interfaces.EnvironmentRepository
}

func NewEnvironmentService(environmentRepository interfaces.EnvironmentRepository) *EnvironmentService {
	return &EnvironmentService{
		environmentRepository: environmentRepository,
	}
}

func (es *EnvironmentService) ValidateUserPermission(ctx context.Context, userId int64, environmentId int64) error {
	users, err := es.environmentRepository.GetEnvironmentAuthorizedUsers(ctx, environmentId)
	if err != nil {
		return err
	}

	found := false
	for _, user := range users {
		if user == userId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not authorized")
	}
	return nil
}

func (es *EnvironmentService) RandomPrefix(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

package application

import (
	"context"
	"database/sql"
	"errors"
	"time"

	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
)

const runnerRegistrationTokenTTL = time.Hour

type RunnerApplication struct {
	runnerRepository    interfaces.RunnerRepository
	organizationService *service.OrganizationService
	tokenService        *service.TokenService
}

func NewRunnerApplication(
	runnerRepository interfaces.RunnerRepository,
	organizationService *service.OrganizationService,
	tokenService *service.TokenService,
) *RunnerApplication {
	return &RunnerApplication{
		runnerRepository:    runnerRepository,
		organizationService: organizationService,
		tokenService:        tokenService,
	}
}

func (ra *RunnerApplication) CreateRunner(
	ctx context.Context,
	organizationId int64,
	userId int64,
) (*value.CreateRunnerResult, error) {
	err := ra.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	token, err := ra.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(runnerRegistrationTokenTTL)

	runner, err := ra.runnerRepository.CreateRunnerWithRegistrationToken(
		ctx,
		organizationId,
		ra.tokenService.HashToken(token),
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &value.CreateRunnerResult{
		Id:        runner.Id,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (ra *RunnerApplication) RegisterRunner(
	ctx context.Context,
	token string,
) error {
	err := ra.runnerRepository.UseRunnerRegistrationToken(ctx, ra.tokenService.HashToken(token))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value.ErrInvalidRunnerRegistrationToken
		}
		return err
	}

	return nil
}

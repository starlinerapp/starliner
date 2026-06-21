package application

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	coreValue "starliner.app/internal/core/domain/value"
)

const runnerRegistrationTokenTTL = time.Hour

type RunnerApplication struct {
	runnerRepository    interfaces.RunnerRepository
	organizationService *service.OrganizationService
	tokenService        *service.TokenService
	queue               port.Queue
}

func NewRunnerApplication(
	runnerRepository interfaces.RunnerRepository,
	organizationService *service.OrganizationService,
	tokenService *service.TokenService,
	queue port.Queue,
) *RunnerApplication {
	return &RunnerApplication{
		runnerRepository:    runnerRepository,
		organizationService: organizationService,
		tokenService:        tokenService,
		queue:               queue,
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
		&organizationId,
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
	name string,
	labels []string,
	maxConcurrentJobs int32,
) error {
	err := ra.runnerRepository.RegisterRunner(
		ctx,
		ra.tokenService.HashToken(token),
		name,
		labels,
		maxConcurrentJobs,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value.ErrInvalidRunnerRegistrationToken
		}
		return err
	}

	return nil
}

func (ra *RunnerApplication) GetOrganizationRunners(
	ctx context.Context,
	organizationId int64,
	userId int64,
) ([]*value.Runner, error) {
	err := ra.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	runners, err := ra.runnerRepository.GetOrganizationRunners(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	globalRunners, err := ra.runnerRepository.GetGlobalRunners(ctx)
	if err != nil {
		return nil, err
	}

	allRunners := append(runners, globalRunners...)

	return value.NewRunners(allRunners), nil
}

func (ra *RunnerApplication) ResolveRunnerId(ctx context.Context, token string) (int64, *int64, error) {
	runnerId, organizationId, err := ra.runnerRepository.ResolveRunnerByRegistrationToken(
		ctx,
		ra.tokenService.HashToken(token),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil, value.ErrInvalidRunnerRegistrationToken
		}
		return 0, nil, err
	}

	return runnerId, organizationId, nil
}

func (ra *RunnerApplication) CreateGlobalRunner(ctx context.Context) (*value.CreateRunnerResult, error) {
	token, err := ra.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(runnerRegistrationTokenTTL)

	runner, err := ra.runnerRepository.CreateRunnerWithRegistrationToken(
		ctx,
		nil,
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

func (ra *RunnerApplication) ListGlobalRunners(ctx context.Context) ([]*value.Runner, error) {
	runners, err := ra.runnerRepository.GetGlobalRunners(ctx)
	if err != nil {
		return nil, err
	}

	return value.NewRunners(runners), nil
}

func (ra *RunnerApplication) DeleteGlobalRunner(ctx context.Context, runnerId int64) error {
	runner, err := ra.runnerRepository.GetRunnerById(ctx, runnerId)
	if err != nil {
		return err
	}

	if runner.OrganizationId != nil {
		return sql.ErrNoRows
	}

	if err := ra.runnerRepository.DeleteGlobalRunner(ctx, runnerId); err != nil {
		return err
	}

	err = ra.queue.PublishRunnerDeleted(&coreValue.RunnerDeleted{RunnerId: runnerId})
	if err != nil {
		log.Printf("error publishing runner deleted: %v", err)
	}

	return nil
}

func (ra *RunnerApplication) HandleRunnerStatusChanged(
	ctx context.Context,
	status *coreValue.RunnerStatusChanged,
) error {
	return ra.runnerRepository.UpdateRunnerStatus(ctx, status.RunnerId, string(status.Status))
}

func (ra *RunnerApplication) DeleteRunner(
	ctx context.Context,
	organizationId int64,
	runnerId int64,
	userId int64,
) error {
	err := ra.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	_, err = ra.runnerRepository.GetRunnerByOrganization(ctx, runnerId, organizationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	if err := ra.runnerRepository.DeleteRunner(ctx, runnerId, organizationId); err != nil {
		return err
	}

	err = ra.queue.PublishRunnerDeleted(&coreValue.RunnerDeleted{RunnerId: runnerId})
	if err != nil {
		log.Printf("error publishing runner deleted: %v", err)
	}

	return nil
}

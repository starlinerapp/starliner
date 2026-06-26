package application

import (
	"context"
	"log"
	"time"

	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/builder/domain/value"
	coreValue "starliner.app/internal/core/domain/value"
)

type HeartbeatApplication struct {
	runnerStore       port.RunnerStore
	queue             port.Queue
	buildApplication  *BuildApplication
	LeaseTTL          time.Duration
	HeartbeatInterval time.Duration
}

func NewHeartbeatApplication(
	runnerStore port.RunnerStore,
	queue port.Queue,
	buildApplication *BuildApplication,
) *HeartbeatApplication {
	leaseTTL := 10 * time.Second

	return &HeartbeatApplication{
		runnerStore:       runnerStore,
		queue:             queue,
		buildApplication:  buildApplication,
		LeaseTTL:          leaseTTL,
		HeartbeatInterval: leaseTTL / 2,
	}
}

func (a *HeartbeatApplication) AcknowledgeHeartbeat(
	ctx context.Context,
	runnerId int64,
	organizationId *int64,
	maxConcurrentJobs int32,
	activeJobs int32,
) (time.Duration, error) {
	if err := validateRunnerCapacity(maxConcurrentJobs, activeJobs); err != nil {
		return 0, err
	}

	deleted, err := a.runnerStore.IsRunnerDeleted(ctx, runnerId)
	if err != nil {
		return 0, err
	}
	if deleted {
		return 0, value.ErrRunnerDeleted
	}

	state, err := a.runnerStore.GetRunner(ctx, runnerId)
	if err != nil {
		return 0, err
	}

	status := coreValue.RunnerStatusOffline
	if state != nil && state.Status != "" {
		status = state.Status
	}

	nextState := port.RunnerRuntimeState{
		OrganizationId:    organizationId,
		Status:            coreValue.RunnerStatusOnline,
		MaxConcurrentJobs: maxConcurrentJobs,
		ActiveJobs:        activeJobs,
	}

	if err := a.runnerStore.UpsertHeartbeat(ctx, runnerId, nextState, a.LeaseTTL); err != nil {
		return 0, err
	}

	if status != coreValue.RunnerStatusOnline {
		if err := a.queue.PublishRunnerStatusChanged(&coreValue.RunnerStatusChanged{
			RunnerId: runnerId,
			Status:   coreValue.RunnerStatusOnline,
		}); err != nil {
			return 0, err
		}
	}

	return a.HeartbeatInterval, nil
}

func (a *HeartbeatApplication) CheckMissedHeartbeats(ctx context.Context) error {
	runnerIds, err := a.runnerStore.ListMonitoredRunners(ctx)
	if err != nil {
		return err
	}

	for _, runnerId := range runnerIds {
		state, err := a.runnerStore.GetRunner(ctx, runnerId)
		if err != nil {
			log.Printf("failed to get runner %d status: %v", runnerId, err)
			continue
		}
		if state == nil || state.Status != coreValue.RunnerStatusOnline {
			continue
		}

		alive, err := a.runnerStore.IsRunnerAlive(ctx, runnerId)
		if err != nil {
			log.Printf("failed to check runner %d alive status: %v", runnerId, err)
			continue
		}
		if alive {
			continue
		}

		if err := a.runnerStore.SetRunnerStatus(ctx, runnerId, coreValue.RunnerStatusOffline); err != nil {
			log.Printf("failed to set runner %d offline: %v", runnerId, err)
			continue
		}

		if err := a.queue.PublishRunnerStatusChanged(&coreValue.RunnerStatusChanged{
			RunnerId: runnerId,
			Status:   coreValue.RunnerStatusOffline,
		}); err != nil {
			log.Printf("failed to publish runner %d status change: %v", runnerId, err)
			continue
		}

		a.buildApplication.HandleRunnerStatusChanged(&coreValue.RunnerStatusChanged{
			RunnerId: runnerId,
			Status:   coreValue.RunnerStatusOffline,
		})
	}

	return nil
}

func validateRunnerCapacity(maxConcurrentJobs int32, activeJobs int32) error {
	if maxConcurrentJobs <= 0 {
		return value.ErrInvalidRunnerCapacity
	}
	if activeJobs < 0 {
		return value.ErrInvalidRunnerCapacity
	}
	if activeJobs > maxConcurrentJobs {
		return value.ErrInvalidRunnerCapacity
	}

	return nil
}

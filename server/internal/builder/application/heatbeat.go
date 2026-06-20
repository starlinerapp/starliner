package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"starliner.app/internal/builder/domain/port"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

var ErrRunnerDeleted = errors.New("runner deleted")

type HeartbeatApplication struct {
	livenessStore     corePort.LivenessStore
	queue             port.Queue
	LeaseTTL          time.Duration
	HeartbeatInterval time.Duration
}

func NewHeartbeatApplication(
	livenessStore corePort.LivenessStore,
	queue port.Queue,
) *HeartbeatApplication {
	leaseTTL := 10 * time.Second

	return &HeartbeatApplication{
		livenessStore:     livenessStore,
		queue:             queue,
		LeaseTTL:          leaseTTL,
		HeartbeatInterval: leaseTTL / 2,
	}
}

func (a *HeartbeatApplication) AcknowledgeHeartbeat(ctx context.Context, runnerId int64) (time.Duration, error) {
	deleted, err := a.livenessStore.IsRunnerDeleted(ctx, runnerId)
	if err != nil {
		return 0, err
	}
	if deleted {
		return 0, ErrRunnerDeleted
	}

	key := fmt.Sprintf("runner:%d", runnerId)

	status, err := a.livenessStore.GetRunnerStatus(ctx, runnerId)
	if err != nil {
		return 0, err
	}

	if err := a.livenessStore.MarkAlive(ctx, key, a.LeaseTTL); err != nil {
		return 0, err
	}

	if err := a.livenessStore.AddMonitoredRunner(ctx, runnerId); err != nil {
		return 0, err
	}

	if status != string(value.RunnerStatusOnline) {
		if err := a.livenessStore.SetRunnerStatus(ctx, runnerId, string(value.RunnerStatusOnline)); err != nil {
			return 0, err
		}

		if err := a.queue.PublishRunnerStatusChanged(&value.RunnerStatusChanged{
			RunnerId: runnerId,
			Status:   value.RunnerStatusOnline,
		}); err != nil {
			return 0, err
		}
	}

	return a.HeartbeatInterval, nil
}

func (a *HeartbeatApplication) CheckMissedHeartbeats(ctx context.Context) error {
	runnerIds, err := a.livenessStore.ListMonitoredRunners(ctx)
	if err != nil {
		return err
	}

	for _, runnerId := range runnerIds {
		status, err := a.livenessStore.GetRunnerStatus(ctx, runnerId)
		if err != nil {
			return err
		}

		if status != string(value.RunnerStatusOnline) {
			continue
		}

		alive, err := a.livenessStore.IsAlive(ctx, fmt.Sprintf("runner:%d", runnerId))
		if err != nil {
			return err
		}

		if alive {
			continue
		}

		if err := a.livenessStore.SetRunnerStatus(ctx, runnerId, string(value.RunnerStatusOffline)); err != nil {
			return err
		}

		if err := a.queue.PublishRunnerStatusChanged(&value.RunnerStatusChanged{
			RunnerId: runnerId,
			Status:   value.RunnerStatusOffline,
		}); err != nil {
			return err
		}
	}

	return nil
}

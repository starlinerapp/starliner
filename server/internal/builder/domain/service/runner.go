package service

import (
	"context"
	"errors"

	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

var ErrNoEligibleRunner = errors.New("no eligible runner")

type RunnerService struct {
	runnerStore corePort.RunnerStore
}

func NewRunnerService(runnerStore corePort.RunnerStore) *RunnerService {
	return &RunnerService{
		runnerStore: runnerStore,
	}
}

func (s *RunnerService) SelectRunner(ctx context.Context, organizationId int64) (int64, error) {
	runnerIds, err := s.runnerStore.ListOrgRunners(ctx, organizationId)
	if err != nil {
		return 0, err
	}

	var (
		selectedID   int64
		selectedLoad float64
		found        bool
	)

	for _, runnerId := range runnerIds {
		state, ok, err := s.eligibleRunner(ctx, runnerId)
		if err != nil {
			return 0, err
		}
		if !ok {
			continue
		}

		load := runnerLoad(state)
		if !found || load < selectedLoad || (load == selectedLoad && runnerId < selectedID) {
			selectedID = runnerId
			selectedLoad = load
			found = true
		}
	}

	if !found {
		return 0, ErrNoEligibleRunner
	}

	return selectedID, nil
}

func (s *RunnerService) eligibleRunner(
	ctx context.Context,
	runnerId int64,
) (*corePort.RunnerRuntimeState, bool, error) {
	deleted, err := s.runnerStore.IsRunnerDeleted(ctx, runnerId)
	if err != nil {
		return nil, false, err
	}
	if deleted {
		return nil, false, nil
	}

	alive, err := s.runnerStore.IsRunnerAlive(ctx, runnerId)
	if err != nil {
		return nil, false, err
	}
	if !alive {
		return nil, false, nil
	}

	state, err := s.runnerStore.GetRunner(ctx, runnerId)
	if err != nil {
		return nil, false, err
	}
	if state == nil || state.Status != value.RunnerStatusOnline {
		return nil, false, nil
	}
	if state.ActiveJobs >= state.MaxConcurrentJobs {
		return nil, false, nil
	}

	return state, true, nil
}

func runnerLoad(state *corePort.RunnerRuntimeState) float64 {
	return float64(state.ActiveJobs) / float64(state.MaxConcurrentJobs)
}

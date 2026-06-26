package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
)

const (
	monitoredRunnersKey   = "runners:monitored"
	runnersGlobalKey      = "runners:global"
	runnerKeyFmt          = "runner:%d"
	runnerLiveKeyFmt      = "live:runner:%d"
	runnerDeletedKeyFmt   = "runner:%d:deleted"
	runnersOrgKeyFmt      = "runners:org:%d"
	runnerFieldOrgID      = "org_id"
	runnerFieldStatus     = "status"
	runnerFieldMaxJobs    = "max_jobs"
	runnerFieldActiveJobs = "active_jobs"
)

func (s *Store) UpsertHeartbeat(
	ctx context.Context,
	runnerId int64,
	state port.RunnerRuntimeState,
	leaseTTL time.Duration,
) error {
	runnerKey := fmt.Sprintf(runnerKeyFmt, runnerId)
	liveKey := fmt.Sprintf(runnerLiveKeyFmt, runnerId)

	pipe := s.client.Pipeline()
	pipe.HSet(ctx, runnerKey, map[string]any{
		runnerFieldOrgID:      organizationIDToRedis(state.OrganizationId),
		runnerFieldStatus:     string(state.Status),
		runnerFieldMaxJobs:    state.MaxConcurrentJobs,
		runnerFieldActiveJobs: state.ActiveJobs,
	})
	pipe.Set(ctx, liveKey, "1", leaseTTL)
	pipe.SAdd(ctx, monitoredRunnersKey, runnerId)
	if state.OrganizationId != nil {
		pipe.SAdd(ctx, fmt.Sprintf(runnersOrgKeyFmt, *state.OrganizationId), runnerId)
	} else {
		pipe.SAdd(ctx, runnersGlobalKey, runnerId)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Store) GetRunner(ctx context.Context, runnerId int64) (*port.RunnerRuntimeState, error) {
	values, err := s.client.HGetAll(ctx, fmt.Sprintf(runnerKeyFmt, runnerId)).Result()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}

	orgID, err := organizationIDFromRedis(values[runnerFieldOrgID])
	if err != nil {
		return nil, fmt.Errorf("parse org_id: %w", err)
	}

	maxJobs, err := strconv.ParseInt(values[runnerFieldMaxJobs], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parse max_jobs: %w", err)
	}

	activeJobs, err := strconv.ParseInt(values[runnerFieldActiveJobs], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parse active_jobs: %w", err)
	}

	return &port.RunnerRuntimeState{
		OrganizationId:    orgID,
		Status:            value.RunnerStatus(values[runnerFieldStatus]),
		MaxConcurrentJobs: int32(maxJobs),
		ActiveJobs:        int32(activeJobs),
	}, nil
}

func (s *Store) IsRunnerDeleted(ctx context.Context, runnerId int64) (bool, error) {
	count, err := s.client.Exists(ctx, fmt.Sprintf(runnerDeletedKeyFmt, runnerId)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Store) IsRunnerAlive(ctx context.Context, runnerId int64) (bool, error) {
	count, err := s.client.Exists(ctx, fmt.Sprintf(runnerLiveKeyFmt, runnerId)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Store) ListOrgRunners(ctx context.Context, orgId int64) ([]int64, error) {
	return parseRunnerSet(s.client.SMembers(ctx, fmt.Sprintf(runnersOrgKeyFmt, orgId)).Result())
}

func (s *Store) ListGlobalRunners(ctx context.Context) ([]int64, error) {
	return parseRunnerSet(s.client.SMembers(ctx, runnersGlobalKey).Result())
}

func (s *Store) ListMonitoredRunners(ctx context.Context) ([]int64, error) {
	return parseRunnerSet(s.client.SMembers(ctx, monitoredRunnersKey).Result())
}

func (s *Store) SetRunnerStatus(ctx context.Context, runnerId int64, status value.RunnerStatus) error {
	return s.client.HSet(ctx, fmt.Sprintf(runnerKeyFmt, runnerId), runnerFieldStatus, string(status)).Err()
}

func (s *Store) MarkDeleted(ctx context.Context, runnerId int64, state *port.RunnerRuntimeState) error {
	var orgId *int64
	if state != nil {
		orgId = state.OrganizationId
	}

	pipe := s.client.Pipeline()
	if orgId != nil {
		pipe.SRem(ctx, fmt.Sprintf(runnersOrgKeyFmt, *orgId), runnerId)
	} else {
		pipe.SRem(ctx, runnersGlobalKey, runnerId)
	}
	pipe.SRem(ctx, monitoredRunnersKey, runnerId)
	pipe.Del(ctx, fmt.Sprintf(runnerKeyFmt, runnerId))
	pipe.Del(ctx, fmt.Sprintf(runnerLiveKeyFmt, runnerId))
	pipe.Set(ctx, fmt.Sprintf(runnerDeletedKeyFmt, runnerId), "1", 0)
	_, err := pipe.Exec(ctx)
	return err
}

func organizationIDToRedis(organizationId *int64) string {
	if organizationId == nil {
		return ""
	}

	return strconv.FormatInt(*organizationId, 10)
}

func organizationIDFromRedis(raw string) (*int64, error) {
	if raw == "" {
		return nil, nil
	}

	orgID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, err
	}

	return &orgID, nil
}

func parseRunnerSet(members []string, err error) ([]int64, error) {
	if err != nil {
		return nil, err
	}

	runnerIds := make([]int64, 0, len(members))
	for _, member := range members {
		runnerId, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			continue
		}
		runnerIds = append(runnerIds, runnerId)
	}

	return runnerIds, nil
}

var _ port.RunnerStore = (*Store)(nil)

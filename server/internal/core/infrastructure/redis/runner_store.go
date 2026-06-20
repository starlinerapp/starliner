package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

const (
	monitoredRunnersKey   = "runners:monitored"
	runnerKeyFmt          = "runner:%d"
	runnerLiveKeyFmt      = "live:runner:%d"
	runnerDeletedKeyFmt   = "runner:%d:deleted"
	runnersOrgKeyFmt      = "runners:org:%d"
	runnerFieldOrgID      = "org_id"
	runnerFieldStatus     = "status"
	runnerFieldMaxJobs    = "max_jobs"
	runnerFieldActiveJobs = "active_jobs"
)

func (c *Client) UpsertHeartbeat(
	ctx context.Context,
	runnerId int64,
	state port.RunnerRuntimeState,
	leaseTTL time.Duration,
) error {
	runnerKey := fmt.Sprintf(runnerKeyFmt, runnerId)
	liveKey := fmt.Sprintf(runnerLiveKeyFmt, runnerId)

	pipe := c.client.Pipeline()
	pipe.HSet(ctx, runnerKey, map[string]any{
		runnerFieldOrgID:      state.OrganizationId,
		runnerFieldStatus:     string(state.Status),
		runnerFieldMaxJobs:    state.MaxConcurrentJobs,
		runnerFieldActiveJobs: state.ActiveJobs,
	})
	pipe.Set(ctx, liveKey, "1", leaseTTL)
	pipe.SAdd(ctx, monitoredRunnersKey, runnerId)
	if state.OrganizationId > 0 {
		pipe.SAdd(ctx, fmt.Sprintf(runnersOrgKeyFmt, state.OrganizationId), runnerId)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c *Client) GetRunner(ctx context.Context, runnerId int64) (*port.RunnerRuntimeState, error) {
	values, err := c.client.HGetAll(ctx, fmt.Sprintf(runnerKeyFmt, runnerId)).Result()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}

	orgID, err := strconv.ParseInt(values[runnerFieldOrgID], 10, 64)
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

func (c *Client) IsRunnerDeleted(ctx context.Context, runnerId int64) (bool, error) {
	count, err := c.client.Exists(ctx, fmt.Sprintf(runnerDeletedKeyFmt, runnerId)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Client) IsRunnerAlive(ctx context.Context, runnerId int64) (bool, error) {
	count, err := c.client.Exists(ctx, fmt.Sprintf(runnerLiveKeyFmt, runnerId)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Client) ListOrgRunners(ctx context.Context, orgId int64) ([]int64, error) {
	return parseRunnerSet(c.client.SMembers(ctx, fmt.Sprintf(runnersOrgKeyFmt, orgId)).Result())
}

func (c *Client) ListMonitoredRunners(ctx context.Context) ([]int64, error) {
	return parseRunnerSet(c.client.SMembers(ctx, monitoredRunnersKey).Result())
}

func (c *Client) SetRunnerStatus(ctx context.Context, runnerId int64, status value.RunnerStatus) error {
	return c.client.HSet(ctx, fmt.Sprintf(runnerKeyFmt, runnerId), runnerFieldStatus, string(status)).Err()
}

func (c *Client) MarkDeleted(ctx context.Context, runnerId int64, orgId int64) error {
	pipe := c.client.Pipeline()
	if orgId > 0 {
		pipe.SRem(ctx, fmt.Sprintf(runnersOrgKeyFmt, orgId), runnerId)
	}
	pipe.SRem(ctx, monitoredRunnersKey, runnerId)
	pipe.Del(ctx, fmt.Sprintf(runnerKeyFmt, runnerId))
	pipe.Del(ctx, fmt.Sprintf(runnerLiveKeyFmt, runnerId))
	pipe.Set(ctx, fmt.Sprintf(runnerDeletedKeyFmt, runnerId), "1", 0)
	_, err := pipe.Exec(ctx)
	return err
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

var _ port.RunnerStore = (*Client)(nil)

package redis

import (
	"context"
	"fmt"
	"strconv"

	"starliner.app/internal/builder/domain/port"
)

const (
	runnerDispatchedBuildsKeyFmt = "runner:%d:dispatched_builds"
	buildDispatchKeyFmt          = "build:%d:dispatch"
	buildDispatchFieldRunnerID   = "runner_id"
	buildDispatchFieldDeployment = "deployment_id"
)

func (s *Store) Register(
	ctx context.Context,
	runnerId int64,
	buildId int64,
	deploymentId int64,
) error {
	runnerKey := fmt.Sprintf(runnerDispatchedBuildsKeyFmt, runnerId)
	buildKey := fmt.Sprintf(buildDispatchKeyFmt, buildId)

	pipe := s.client.Pipeline()
	pipe.SAdd(ctx, runnerKey, buildId)
	pipe.HSet(ctx, buildKey, map[string]any{
		buildDispatchFieldRunnerID:   runnerId,
		buildDispatchFieldDeployment: deploymentId,
	})
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Store) Unregister(ctx context.Context, buildId int64) error {
	buildKey := fmt.Sprintf(buildDispatchKeyFmt, buildId)

	values, err := s.client.HGetAll(ctx, buildKey).Result()
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}

	runnerId, err := strconv.ParseInt(values[buildDispatchFieldRunnerID], 10, 64)
	if err != nil {
		return fmt.Errorf("parse runner_id: %w", err)
	}

	pipe := s.client.Pipeline()
	pipe.SRem(ctx, fmt.Sprintf(runnerDispatchedBuildsKeyFmt, runnerId), buildId)
	pipe.Del(ctx, buildKey)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Store) ListByRunner(ctx context.Context, runnerId int64) ([]port.DispatchedBuild, error) {
	members, err := s.client.SMembers(ctx, fmt.Sprintf(runnerDispatchedBuildsKeyFmt, runnerId)).Result()
	if err != nil {
		return nil, err
	}

	builds := make([]port.DispatchedBuild, 0, len(members))
	for _, member := range members {
		buildId, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			continue
		}

		values, err := s.client.HGetAll(ctx, fmt.Sprintf(buildDispatchKeyFmt, buildId)).Result()
		if err != nil {
			return nil, err
		}
		if len(values) == 0 {
			continue
		}

		deploymentId, err := strconv.ParseInt(values[buildDispatchFieldDeployment], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse deployment_id for build %d: %w", buildId, err)
		}

		builds = append(builds, port.DispatchedBuild{
			BuildId:      buildId,
			DeploymentId: deploymentId,
		})
	}

	return builds, nil
}

var _ port.BuildDispatchStore = (*Store)(nil)

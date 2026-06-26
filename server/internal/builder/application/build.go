package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/builder/domain/service"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

const (
	buildDispatchLockTTL = 24 * time.Hour
	buildCompleteLockTTL = 24 * time.Hour
	runnerOfflineMessage = "build failed: runner went offline"
)

type BuildApplication struct {
	cfg                *conf.Config
	queue              port.Queue
	logPublisher       port.LogPublisher
	runnerService      *service.RunnerService
	lease              corePort.Lease
	buildDispatchStore port.BuildDispatchStore
}

func NewBuildApplication(
	cfg *conf.Config,
	queue port.Queue,
	logPublisher port.LogPublisher,
	runnerService *service.RunnerService,
	lease corePort.Lease,
	buildDispatchStore port.BuildDispatchStore,
) *BuildApplication {
	return &BuildApplication{
		cfg:                cfg,
		queue:              queue,
		logPublisher:       logPublisher,
		runnerService:      runnerService,
		lease:              lease,
		buildDispatchStore: buildDispatchStore,
	}
}

func (ba *BuildApplication) HandleBuildTriggered(build *value.TriggerBuild) {
	ctx := context.Background()

	publishLogLine := func(line string) {
		if ba.logPublisher == nil {
			return
		}
		if err := ba.logPublisher.PublishLogChunk(build.BuildId, []byte(line)); err != nil {
			log.Printf("failed to publish log chunk: %v", err)
		}
	}

	publishCompleted := func(commitHash, tag *string, imageName *string, logs string, status value.BuildStatus) {
		if err := ba.queue.PublishBuildCompleted(&value.BuildCompleted{
			BuildId:          build.BuildId,
			DeploymentId:     build.DeploymentId,
			ImageRegistryUrl: ba.cfg.ImageRegistryUrl,
			ImageName:        imageName,
			CommitHash:       commitHash,
			Tag:              commitHash,
			Logs:             logs,
			BuildStatus:      status,
		}); err != nil {
			log.Printf("failed to publish: %v", err)
		}
	}

	runnerId, err := ba.runnerService.SelectRunner(ctx, build.OrganizationId)
	if err != nil {
		msg := "no eligible runner available for organization"
		if !errors.Is(err, service.ErrNoEligibleRunner) {
			msg = fmt.Sprintf("select runner: %v", err)
		}
		publishLogLine(msg + "\n")
		if ba.logPublisher != nil {
			_ = ba.logPublisher.PublishLogEnd(build.BuildId)
		}
		publishCompleted(nil, nil, nil, msg, value.BuildStatusFailed)
		return
	}

	acquired, err := ba.lease.TryLease(
		ctx,
		fmt.Sprintf("build:dispatch:%d", build.BuildId),
		buildDispatchLockTTL,
	)
	if err != nil {
		msg := fmt.Sprintf("claim build dispatch lock: %v", err)
		publishLogLine(msg + "\n")
		if ba.logPublisher != nil {
			_ = ba.logPublisher.PublishLogEnd(build.BuildId)
		}
		publishCompleted(nil, nil, nil, msg, value.BuildStatusFailed)
		return
	}
	if !acquired {
		log.Printf("build %d already dispatched, skipping duplicate trigger", build.BuildId)
		return
	}

	job := &value.RunnerBuildJob{
		RunnerId:          runnerId,
		BuildId:           build.BuildId,
		DeploymentId:      build.DeploymentId,
		ImageName:         build.ImageName,
		ImageRegistryUrl:  ba.cfg.ImageRegistryUrl,
		GitUrl:            build.GitUrl,
		BranchName:        build.BranchName,
		AccessToken:       build.AccessToken,
		RegistryPushToken: build.RegistryPushToken,
		RootDirectory:     build.RootDirectory,
		DockerfilePath:    build.DockerfilePath,
		Args:              build.Args,
	}

	if err := ba.queue.PublishRunnerJob(job); err != nil {
		msg := fmt.Sprintf("enqueue build for runner: %v", err)
		publishLogLine(msg + "\n")
		if ba.logPublisher != nil {
			_ = ba.logPublisher.PublishLogEnd(build.BuildId)
		}
		publishCompleted(nil, nil, nil, msg, value.BuildStatusFailed)
		return
	}

	if err := ba.buildDispatchStore.Register(ctx, runnerId, build.BuildId, build.DeploymentId); err != nil {
		log.Printf("failed to register dispatched build %d for runner %d: %v", build.BuildId, runnerId, err)
	}

	log.Printf("enqueued build %d for runner %d", build.BuildId, runnerId)
}

func (ba *BuildApplication) HandleRunnerJobResult(result *value.RunnerBuildResult) {
	if result == nil {
		return
	}

	ctx := context.Background()
	if !ba.tryMarkBuildComplete(ctx, result.BuildId) {
		return
	}

	if err := ba.buildDispatchStore.Unregister(ctx, result.BuildId); err != nil {
		log.Printf("failed to unregister dispatched build %d: %v", result.BuildId, err)
	}

	if ba.logPublisher != nil {
		if err := ba.logPublisher.PublishLogEnd(result.BuildId); err != nil {
			log.Printf("failed to publish log end: %v", err)
		}
	}

	var commitHash, imageName *string
	if result.CommitHash != "" {
		commitHash = &result.CommitHash
	}
	if result.ImageName != "" {
		imageName = &result.ImageName
	}

	if err := ba.queue.PublishBuildCompleted(&value.BuildCompleted{
		BuildId:          result.BuildId,
		DeploymentId:     result.DeploymentId,
		ImageRegistryUrl: ba.cfg.ImageRegistryUrl,
		ImageName:        imageName,
		CommitHash:       commitHash,
		Tag:              commitHash,
		Logs:             result.Logs,
		BuildStatus:      result.Status,
	}); err != nil {
		log.Printf("failed to publish build completed: %v", err)
	}
}

func (ba *BuildApplication) HandleRunnerStatusChanged(status *value.RunnerStatusChanged) {
	if status == nil || status.Status != value.RunnerStatusOffline {
		return
	}

	ctx := context.Background()
	builds, err := ba.buildDispatchStore.ListByRunner(ctx, status.RunnerId)
	if err != nil {
		log.Printf("failed to list dispatched builds for runner %d: %v", status.RunnerId, err)
		return
	}

	for _, build := range builds {
		ba.failBuildForRunnerOffline(ctx, build.BuildId, build.DeploymentId)
	}
}

func (ba *BuildApplication) failBuildForRunnerOffline(ctx context.Context, buildId int64, deploymentId int64) {
	if !ba.tryMarkBuildComplete(ctx, buildId) {
		return
	}

	if err := ba.buildDispatchStore.Unregister(ctx, buildId); err != nil {
		log.Printf("failed to unregister dispatched build %d: %v", buildId, err)
	}

	if ba.logPublisher != nil {
		if err := ba.logPublisher.PublishLogChunk(buildId, []byte(runnerOfflineMessage+"\n")); err != nil {
			log.Printf("failed to publish log chunk: %v", err)
		}
		if err := ba.logPublisher.PublishLogEnd(buildId); err != nil {
			log.Printf("failed to publish log end: %v", err)
		}
	}

	if err := ba.queue.PublishBuildCompleted(&value.BuildCompleted{
		BuildId:          buildId,
		DeploymentId:     deploymentId,
		ImageRegistryUrl: ba.cfg.ImageRegistryUrl,
		BuildStatus:      value.BuildStatusFailed,
		Logs:             runnerOfflineMessage,
	}); err != nil {
		log.Printf("failed to publish build completed: %v", err)
	}
}

func (ba *BuildApplication) tryMarkBuildComplete(ctx context.Context, buildId int64) bool {
	acquired, err := ba.lease.TryLease(
		ctx,
		fmt.Sprintf("build:complete:%d", buildId),
		buildCompleteLockTTL,
	)
	if err != nil {
		log.Printf("failed to claim build complete lock for build %d: %v", buildId, err)
		return false
	}

	return acquired
}

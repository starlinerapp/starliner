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

const buildDispatchLockTTL = 24 * time.Hour

type BuildApplication struct {
	cfg             *conf.Config
	queue           port.Queue
	logPublisher    port.LogPublisher
	runnerService   *service.RunnerService
	dispatchLimiter corePort.AcquireLimiter
}

func NewBuildApplication(
	cfg *conf.Config,
	queue port.Queue,
	logPublisher port.LogPublisher,
	runnerService *service.RunnerService,
	dispatchLimiter corePort.AcquireLimiter,
) *BuildApplication {
	return &BuildApplication{
		cfg:             cfg,
		queue:           queue,
		logPublisher:    logPublisher,
		runnerService:   runnerService,
		dispatchLimiter: dispatchLimiter,
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

	acquired, err := ba.dispatchLimiter.TryAcquire(
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

	log.Printf("enqueued build %d for runner %d", build.BuildId, runnerId)
}

func (ba *BuildApplication) HandleRunnerJobResult(result *value.RunnerBuildResult) {
	if result == nil {
		return
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

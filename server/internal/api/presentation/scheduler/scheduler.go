package scheduler

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/api/application"
	"starliner.app/internal/core/util/concurrent"
	"time"
)

type Scheduler struct {
	deploymentApplication *application.DeploymentApplication
}

func RegisterScheduler(lc fx.Lifecycle, s *Scheduler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Start()
		},
	})
}

func NewScheduler(
	deploymentApplication *application.DeploymentApplication,
) *Scheduler {
	return &Scheduler{
		deploymentApplication: deploymentApplication,
	}
}

func (s *Scheduler) Start() error {
	go concurrent.WithRecovery(context.Background(), "RequestDeploymentStatus", func() error {
		interval := 2 * time.Second

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := s.deploymentApplication.RequestDeploymentStatus(); err != nil {
				log.Printf("failed to request deployment status: %v", err)
			}
		}

		return nil
	})
	return nil
}

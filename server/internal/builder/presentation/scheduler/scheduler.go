package scheduler

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/core/util/concurrent"
)

type Scheduler struct {
	heartbeatApplication *application.HeartbeatApplication
}

func RegisterScheduler(lc fx.Lifecycle, s *Scheduler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Start()
		},
	})
}

func NewScheduler(
	heartbeatApplication *application.HeartbeatApplication,
) *Scheduler {
	return &Scheduler{
		heartbeatApplication: heartbeatApplication,
	}
}

func (s *Scheduler) Start() error {
	go concurrent.WithRecovery(context.Background(), "CheckMissedRunnerHeartbeats", func() error {
		interval := s.heartbeatApplication.HeartbeatInterval

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := s.heartbeatApplication.CheckMissedHeartbeats(context.Background()); err != nil {
				log.Printf("failed to check missed runner heartbeats: %v", err)
			}
		}

		return nil
	})
	return nil
}

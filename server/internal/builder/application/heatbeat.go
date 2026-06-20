package application

import (
	"context"
	"fmt"
	"time"

	"starliner.app/internal/core/domain/port"
)

type HeartbeatApplication struct {
	livenessStore port.LivenessStore
	LeaseTTL      time.Duration
}

func NewHeartbeatApplication(
	livenessStore port.LivenessStore,
) *HeartbeatApplication {
	return &HeartbeatApplication{
		livenessStore: livenessStore,
		LeaseTTL:      10 * time.Second,
	}
}

func (a *HeartbeatApplication) AcknowledgeHeartbeat(ctx context.Context, runnerId int64) error {
	return a.livenessStore.MarkAlive(ctx, fmt.Sprintf("runner:%d", runnerId), a.LeaseTTL)
}

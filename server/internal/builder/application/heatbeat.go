package application

import (
	"context"
	"fmt"
	"time"

	"starliner.app/internal/core/domain/port"
)

type HeartbeatApplication struct {
	livenessStore     port.LivenessStore
	LeaseTTL          time.Duration
	HeartbeatInterval time.Duration
}

func NewHeartbeatApplication(
	livenessStore port.LivenessStore,
) *HeartbeatApplication {
	leaseTTL := 10 * time.Second

	return &HeartbeatApplication{
		livenessStore:     livenessStore,
		LeaseTTL:          leaseTTL,
		HeartbeatInterval: leaseTTL / 2,
	}
}

func (a *HeartbeatApplication) AcknowledgeHeartbeat(ctx context.Context, runnerId int64) (time.Duration, error) {
	err := a.livenessStore.MarkAlive(ctx, fmt.Sprintf("runner:%d", runnerId), a.LeaseTTL)
	return a.HeartbeatInterval, err
}

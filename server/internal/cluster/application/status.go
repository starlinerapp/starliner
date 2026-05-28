package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/infrastructure/k8s"
	coreport "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

const (
	clusterReconcileCooldown = 10 * time.Minute
	clusterReconcileKeyFmt   = "starliner:cluster:reconcile:%d"
)

type StatusApplication struct {
	health         port.Health
	pubsub         port.Pubsub
	acquireLimiter coreport.AcquireLimiter
}

func NewStatusApplication(
	health port.Health,
	pubsub port.Pubsub,
	acquireLimiter coreport.AcquireLimiter,
) *StatusApplication {
	return &StatusApplication{
		health:         health,
		pubsub:         pubsub,
		acquireLimiter: acquireLimiter,
	}
}

func (sa *StatusApplication) HandleRequestDeploymentStatus(d *value.Deployment) {
	releaseName := d.DeploymentName
	health, err := sa.health.CheckPodsHealthy(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to check pods health: %v\n", err)
		if k8s.IsClusterUnreachable(err) && d.ClusterId != 0 && d.ProvisioningId != "" {
			sa.RequestClusterReconcile(d)
		}
		return
	}

	err = sa.pubsub.PublishDeploymentStatusResponse(&value.HealthStatus{
		DeploymentId: d.DeploymentId,
		Health:       value.Health(health.Health),
		Status:       health.Status,
	})
	if err != nil {
		log.Printf("failed to publish deployment status: %v\n", err)
	}
}

func (sa *StatusApplication) RequestClusterReconcile(d *value.Deployment) {
	ctx := context.Background()
	key := fmt.Sprintf(clusterReconcileKeyFmt, d.ClusterId)

	allowed, err := sa.acquireLimiter.TryAcquire(ctx, key, clusterReconcileCooldown)
	if err != nil {
		log.Printf("reconcile cooldown check failed for cluster %d: %v\n", d.ClusterId, err)
		return
	}
	if !allowed {
		return
	}

	err = sa.pubsub.PublishReconcileClusterRequest(&value.ReconcileClusterRequest{
		ClusterId:      d.ClusterId,
		OrganizationId: d.OrganizationId,
		ProvisioningId: d.ProvisioningId,
	})
	if err != nil {
		log.Printf("failed to publish reconcile cluster request: %v\n", err)
	}
}

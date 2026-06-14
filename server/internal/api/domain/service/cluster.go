package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
)

const (
	clusterReconcileCooldown = 10 * time.Minute
	clusterReconcileKeyFmt   = "starliner:cluster:reconcile:%d"
)

type ClusterService struct {
	clusterRepository      interfaces.ClusterRepository
	organizationRepository interfaces.OrganizationRepository
	acquireLimiter         corePort.AcquireLimiter
	queue                  port.Queue
	crypto                 corePort.Crypto
}

func NewClusterService(
	clusterRepository interfaces.ClusterRepository,
	organizationRepository interfaces.OrganizationRepository,
	acquireLimiter corePort.AcquireLimiter,
	queue port.Queue,
	crypto corePort.Crypto,
) *ClusterService {
	return &ClusterService{
		clusterRepository:      clusterRepository,
		organizationRepository: organizationRepository,
		acquireLimiter:         acquireLimiter,
		queue:                  queue,
		crypto:                 crypto,
	}
}

func (cs *ClusterService) ValidateClusterReady(ctx context.Context, clusterId int64) error {
	cluster, err := cs.clusterRepository.GetCluster(ctx, clusterId)
	if err != nil {
		return err
	}
	if cluster.Status != entity.ClusterRunning {
		return errors.New("cluster is not ready")
	}
	return nil
}

func (cs *ClusterService) ReconcileIfUnreachable(
	ctx context.Context,
	err error,
	deployment *coreValue.Deployment,
) {
	if !coreValue.IsClusterUnreachable(err) {
		return
	}

	cs.requestReconcile(ctx, &coreValue.ReconcileClusterRequest{
		ClusterId:      deployment.ClusterId,
		OrganizationId: deployment.OrganizationId,
		ProvisioningId: deployment.ProvisioningId,
	})
}

func (cs *ClusterService) requestReconcile(ctx context.Context, req *coreValue.ReconcileClusterRequest) {
	if req.ClusterId == 0 || req.ProvisioningId == "" {
		log.Printf(
			"reconcile not triggered: missing cluster metadata (clusterId=%d provisioningId=%q)\n",
			req.ClusterId,
			req.ProvisioningId,
		)
		return
	}

	cluster, err := cs.clusterRepository.GetCluster(ctx, req.ClusterId)
	if err != nil {
		log.Printf("reconcile skipped: cluster %d: %v\n", req.ClusterId, err)
		return
	}

	if cluster.Status != entity.ClusterRunning {
		log.Printf("reconcile skipped: cluster %d status is %s\n", req.ClusterId, cluster.Status)
		return
	}

	if cluster.ProvisioningId == nil || *cluster.ProvisioningId != req.ProvisioningId {
		log.Printf("reconcile skipped: cluster %d provisioning id mismatch\n", req.ClusterId)
		return
	}

	key := fmt.Sprintf(clusterReconcileKeyFmt, req.ClusterId)
	allowed, err := cs.acquireLimiter.TryAcquire(ctx, key, clusterReconcileCooldown)
	if err != nil {
		log.Printf("reconcile cooldown check failed for cluster %d: %v\n", req.ClusterId, err)
		return
	}
	if !allowed {
		log.Printf("reconcile skipped: cluster %d (cooldown)\n", req.ClusterId)
		return
	}

	credential, err := cs.organizationRepository.GetOrganizationProvisioningCredential(
		ctx,
		req.OrganizationId,
		value.HetznerCredential,
	)
	if err != nil {
		log.Printf("reconcile skipped: cluster %d credential: %v\n", req.ClusterId, err)
		return
	}

	decrypted, err := cs.crypto.Decrypt(credential.Secret)
	if err != nil {
		log.Printf("reconcile skipped: cluster %d decrypt credential: %v\n", req.ClusterId, err)
		return
	}

	err = cs.queue.PublishReconcileCluster(&coreValue.ReconcileCluster{
		Id:                     req.ClusterId,
		ProvisioningId:         req.ProvisioningId,
		ProvisioningCredential: decrypted,
	})
	if err != nil {
		log.Printf("failed to publish reconcile cluster for %d: %v\n", req.ClusterId, err)
		return
	}
	log.Printf("reconcile queued for cluster %d (provisioning id %q)\n", req.ClusterId, req.ProvisioningId)
}

package service

import (
	"context"
	"errors"

	"starliner.app/internal/api/domain/entity"
	_interface "starliner.app/internal/api/domain/repository/interface"
)

type ClusterService struct {
	clusterRepository _interface.ClusterRepository
}

func NewClusterService(
	clusterRepository _interface.ClusterRepository,
) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
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

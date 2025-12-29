package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
	interfaces "starliner.app/pkg/repository/interface"
	"strconv"
)

type ClusterService struct {
	organizationRepository interfaces.OrganizationRepository
	clusterRepository      interfaces.ClusterRepository
	clusterPublisher       *queue.Publisher[*v1.Cluster]
}

func NewClusterService(
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	clusterPublisher *queue.Publisher[*v1.Cluster],
) *ClusterService {
	return &ClusterService{
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		clusterPublisher:       clusterPublisher,
	}
}

func (cs *ClusterService) CreateCluster(name string, organizationId int64) error {
	id := uuid.New().String()
	err := cs.clusterPublisher.Publish(queue.CreateCluster, id, &v1.Cluster{
		Name:           name,
		OrganizationId: organizationId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (cs *ClusterService) DeleteCluster(ctx context.Context, userId int64, clusterId int64) error {
	cluster, err := cs.clusterRepository.GetCluster(ctx, clusterId)
	if err != nil {
		return err
	}

	userOrganizations, err := cs.organizationRepository.GetUserOrganizations(ctx, userId)
	if err != nil {
		return err
	}

	found := false
	for _, org := range userOrganizations {
		if org.Id == cluster.OrganizationId {
			found = true
		}
	}

	if !found {
		return errors.New("organization not found")
	}

	err = cs.clusterPublisher.Publish(queue.DeleteCluster, strconv.FormatInt(clusterId, 10), &v1.Cluster{
		Id:             &clusterId,
		Name:           cluster.Name,
		OrganizationId: cluster.OrganizationId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

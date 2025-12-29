package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"starliner.app/pkg/domain"
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

func (cs *ClusterService) CreateCluster(ctx context.Context, name string, organizationId int64) error {
	cluster, err := cs.clusterRepository.CreateCluster(ctx, name, organizationId, nil, nil, nil)
	if err != nil {
		fmt.Printf("failed to persist cluster in database: %v", err)
	}

	err = cs.clusterPublisher.Publish(queue.CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Name:           name,
		OrganizationId: organizationId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (cs *ClusterService) GetCluster(ctx context.Context, id int64, userId int64) (*domain.Cluster, error) {
	cluster, err := cs.clusterRepository.GetCluster(ctx, id)
	if err != nil {
		return nil, err
	}

	orgs, err := cs.organizationRepository.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, err
	}
	found := false
	for _, org := range orgs {
		if org.Id == cluster.OrganizationId {
			found = true
		}
	}
	if !found {
		return nil, errors.New("organization not found")
	}

	return cluster, nil
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

package service

import (
	"github.com/google/uuid"
	"log"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
)

type ClusterService struct {
	clusterPublisher *queue.Publisher[*v1.Cluster]
}

func NewClusterService(clusterPublisher *queue.Publisher[*v1.Cluster]) *ClusterService {
	return &ClusterService{
		clusterPublisher: clusterPublisher,
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

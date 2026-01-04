package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/infrastructure/queue/proto/v1"
	interfaces "starliner.app/internal/repository/interface"
	"starliner.app/internal/service/model"
	"strconv"
)

type ClusterService struct {
	organizationRepository interfaces.OrganizationRepository
	clusterRepository      interfaces.ClusterRepository
	organizationService    *OrganizationService
	cryptoService          *CryptoService
	clusterPublisher       *queue.Publisher[*v1.Cluster]
}

func NewClusterService(
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	organizationService *OrganizationService,
	cryptoService *CryptoService,
	clusterPublisher *queue.Publisher[*v1.Cluster],
) *ClusterService {
	return &ClusterService{
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		organizationService:    organizationService,
		cryptoService:          cryptoService,
		clusterPublisher:       clusterPublisher,
	}
}

func (cs *ClusterService) CreateCluster(ctx context.Context, userId int64, name string, organizationId int64) (*model.Cluster, error) {
	orgs, err := cs.organizationService.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, err
	}

	found := false
	for _, org := range orgs {
		if org.Id == organizationId {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("not permitted")
	}

	cluster, err := cs.clusterRepository.CreateCluster(ctx, name, organizationId)
	if err != nil {
		fmt.Printf("failed to persist cluster in database: %v", err)
	}

	err = cs.clusterPublisher.Publish(queue.CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:             cluster.Id,
		Name:           name,
		OrganizationId: organizationId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return model.NewCluster(cluster), nil
}

func (cs *ClusterService) GetCluster(ctx context.Context, id int64) (*model.Cluster, error) {
	cluster, err := cs.clusterRepository.GetCluster(ctx, id)
	if err != nil {
		return nil, err
	}

	return model.NewCluster(cluster), nil
}

func (cs *ClusterService) GetUserCluster(ctx context.Context, userId int64, id int64) (*model.Cluster, error) {
	cluster, err := cs.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.NewCluster(cluster), nil
}

func (cs *ClusterService) GetClusterPrivateKey(ctx context.Context, id int64, userId int64) ([]byte, error) {
	cluster, err := cs.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if cluster.PrivateKey == nil {
		return nil, fmt.Errorf("cluster private key is not set")
	}

	// The private key was first base64 encoded and then encrypted
	decryptedPrivateKey, err := cs.cryptoService.Decrypt(*cluster.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(decryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	pemBytes, err := cs.cryptoService.EncodePrivateKeyToPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key to PEM: %v", err)
	}
	return pemBytes, nil
}

func (cs *ClusterService) DeleteCluster(ctx context.Context, userId int64, clusterId int64) error {
	cluster, err := cs.clusterRepository.GetUserCluster(ctx, userId, clusterId)
	if err != nil {
		return err
	}

	err = cs.clusterRepository.UpdateClusterStatus(ctx, clusterId, domain.ClusterDeleted)
	if err != nil {
		return err
	}

	err = cs.clusterPublisher.Publish(queue.DeleteCluster, strconv.FormatInt(clusterId, 10), &v1.Cluster{
		Id:             clusterId,
		Name:           cluster.Name,
		OrganizationId: cluster.OrganizationId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

package application

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/repository/interface"
	coreValue "starliner.app/internal/core/domain/value"
	"strconv"
)

type ClusterApplication struct {
	clusterRepository   _interface.ClusterRepository
	organizationService *service.OrganizationService
	crypto              port.Crypto
	queue               port.Queue
}

func NewClusterApplication(
	clusterRepository _interface.ClusterRepository,
	organizationService *service.OrganizationService,
	crypto port.Crypto,
	queue port.Queue,
) *ClusterApplication {
	return &ClusterApplication{
		clusterRepository:   clusterRepository,
		organizationService: organizationService,
		crypto:              crypto,
		queue:               queue,
	}
}

func (ca *ClusterApplication) CreateCluster(ctx context.Context, userId int64, name string, organizationId int64) (*value.Cluster, error) {
	err := ca.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	cluster, err := ca.clusterRepository.CreateCluster(ctx, name, organizationId)
	if err != nil {
		fmt.Printf("failed to persist cluster in database: %v", err)
	}

	err = ca.queue.PublishCreateCluster(&coreValue.Cluster{
		Id:               cluster.Id,
		Name:             cluster.Name,
		OrganizationName: strconv.FormatInt(cluster.OrganizationId, 10),
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return value.NewCluster(cluster), nil
}

func (ca *ClusterApplication) GetCluster(ctx context.Context, id int64) (*value.Cluster, error) {
	cluster, err := ca.clusterRepository.GetCluster(ctx, id)
	if err != nil {
		return nil, err
	}

	return value.NewCluster(cluster), nil
}

func (ca *ClusterApplication) GetUserCluster(ctx context.Context, userId int64, id int64) (*value.Cluster, error) {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return value.NewCluster(cluster), nil
}

func (ca *ClusterApplication) GetClusterPrivateKey(ctx context.Context, id int64, userId int64) ([]byte, error) {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if cluster.PrivateKey == nil {
		return nil, fmt.Errorf("cluster private key is not set")
	}

	// The private key was first base64 encoded and then encrypted
	decryptedPrivateKey, err := ca.crypto.Decrypt(*cluster.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(decryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	pemBytes, err := ca.crypto.EncodePrivateKeyToPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key to PEM: %v", err)
	}
	return pemBytes, nil
}

func (ca *ClusterApplication) DeleteCluster(ctx context.Context, userId int64, clusterId int64) error {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, clusterId)
	if err != nil {
		return err
	}

	err = ca.clusterRepository.UpdateClusterStatus(ctx, clusterId, entity.ClusterDeleted)
	if err != nil {
		return err
	}

	err = ca.queue.PublishDeleteCluster(&coreValue.Cluster{
		Id:               cluster.Id,
		Name:             cluster.Name,
		OrganizationName: strconv.FormatInt(cluster.OrganizationId, 10),
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

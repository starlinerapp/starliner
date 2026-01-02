package service

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"starliner.app/pkg/api/dto/response"
	"starliner.app/pkg/config"
	"starliner.app/pkg/crypto"
	"starliner.app/pkg/domain"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
	interfaces "starliner.app/pkg/repository/interface"
	"strconv"
)

type ClusterService struct {
	cfg                    *config.Config
	organizationRepository interfaces.OrganizationRepository
	clusterRepository      interfaces.ClusterRepository
	clusterPublisher       *queue.Publisher[*v1.Cluster]
}

func NewClusterService(
	cfg *config.Config,
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	clusterPublisher *queue.Publisher[*v1.Cluster],
) *ClusterService {
	return &ClusterService{
		cfg:                    cfg,
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		clusterPublisher:       clusterPublisher,
	}
}

func (cs *ClusterService) CreateCluster(ctx context.Context, name string, organizationId int64) (*response.Cluster, error) {
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

	return &response.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         cluster.Status,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
	}, nil
}

func (cs *ClusterService) GetCluster(ctx context.Context, id int64) (*response.Cluster, error) {
	cluster, err := cs.clusterRepository.GetCluster(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         cluster.Status,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
	}, nil
}

func (cs *ClusterService) GetUserCluster(ctx context.Context, userId int64, id int64) (*response.Cluster, error) {
	cluster, err := cs.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &response.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         cluster.Status,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
	}, nil
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
	encryptionKey, err := base64.StdEncoding.DecodeString(cs.cfg.EncryptionKeyBase64)
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}

	decryptedPrivateKey, err := crypto.Decrypt(*cluster.PrivateKey, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(decryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	privateKey := ed25519.PrivateKey(keyBytes)

	// Serialize to OpenSSH format
	block, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}

	pemBytes := pem.EncodeToMemory(block)
	if pemBytes == nil {
		return nil, fmt.Errorf("failed to encode private key to PEM")
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

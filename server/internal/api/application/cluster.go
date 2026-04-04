package application

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	"strconv"
)

type ClusterApplication struct {
	clusterRepository      interfaces.ClusterRepository
	organizationRepository interfaces.OrganizationRepository
	organizationService    *service.OrganizationService
	crypto                 corePort.Crypto
	queue                  port.Queue
	grpcClient             port.ProvisionerClient
}

func NewClusterApplication(
	clusterRepository interfaces.ClusterRepository,
	organizationRepository interfaces.OrganizationRepository,
	organizationService *service.OrganizationService,
	crypto corePort.Crypto,
	queue port.Queue,
	grpcClient port.ProvisionerClient,
) *ClusterApplication {
	return &ClusterApplication{
		clusterRepository:      clusterRepository,
		organizationRepository: organizationRepository,
		organizationService:    organizationService,
		crypto:                 crypto,
		queue:                  queue,
		grpcClient:             grpcClient,
	}
}

func (ca *ClusterApplication) CreateCluster(ctx context.Context, userId int64, name string, serverType string, organizationId int64) (*value.Cluster, error) {
	err := ca.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	cluster, err := ca.clusterRepository.CreateCluster(ctx, name, serverType, organizationId)
	if err != nil {
		fmt.Printf("failed to persist cluster in database: %v", err)
	}

	credential, err := ca.organizationRepository.GetOrganizationProvisioningCredential(ctx, organizationId, value.HetznerCredential)
	if err != nil {
		return nil, err
	}

	decrypted, err := ca.crypto.Decrypt(credential.Secret)
	if err != nil {
		return nil, err
	}

	err = ca.queue.PublishCreateCluster(&coreValue.ProvisionCluster{
		Id:                     cluster.Id,
		Name:                   cluster.Name,
		ServerType:             coreValue.ServerType(cluster.ServerType),
		OrganizationName:       strconv.FormatInt(cluster.OrganizationId, 10),
		ProvisioningCredential: decrypted,
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

	credential, err := ca.organizationRepository.GetOrganizationProvisioningCredential(ctx, cluster.OrganizationId, value.HetznerCredential)
	if err != nil {
		return err
	}

	decrypted, err := ca.crypto.Decrypt(credential.Secret)
	if err != nil {
		return err
	}

	err = ca.queue.PublishDeleteCluster(&coreValue.DeleteCluster{
		Id:                     cluster.Id,
		ProvisioningId:         *cluster.ProvisioningId,
		ProvisioningCredential: decrypted,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (ca *ClusterApplication) OpenTTY(
	ctx context.Context,
	userId int64,
	clusterId int64,
	stdin io.Reader,
	stdout io.Writer,
	sizes <-chan port.TerminalSize,
) error {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, clusterId)
	if err != nil {
		return err
	}

	if cluster.PrivateKey == nil {
		return fmt.Errorf("cluster private key is not set")
	}
	if cluster.IPv4Address == nil {
		return fmt.Errorf("cluster ipv4 address is not set")
	}

	// The private key was first base64 encoded and then encrypted
	decryptedPrivateKey, err := ca.crypto.Decrypt(*cluster.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %v", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(decryptedPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %v", err)
	}

	pemBytes, err := ca.crypto.EncodePrivateKeyToPEM(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to encode private key to PEM: %v", err)
	}

	// TODO: Don't hardcode user
	return ca.grpcClient.OpenTTY(ctx, *cluster.IPv4Address, "root", pemBytes, stdin, stdout, sizes)
}

func (ca *ClusterApplication) HandleClusterCreated(c *coreValue.ClusterCreated) {
	ctx := context.Background()
	err := ca.clusterRepository.UpdateClusterPulumiStackId(ctx, c.Id, &c.ProvisioningId)
	if err != nil {
		fmt.Printf("failed to persist provisioning id: %v\n", err)
	}

	err = ca.clusterRepository.UpdateClusterIPv4Address(ctx, c.Id, &c.IPv4Address)
	if err != nil {
		fmt.Printf("Failed to persist cluster ip address: %v\n", err)
	}

	encryptedPrivKeyStr, err := ca.crypto.Encrypt(c.PrivateKey)
	if err != nil {
		log.Printf("failed to encrypt private key: %v\n", err)
	}
	err = ca.clusterRepository.UpdateClusterPublicPrivateKey(ctx, c.Id, &c.PublicKey, &encryptedPrivKeyStr)
	if err != nil {
		fmt.Printf("failed to persist cluster public private key: %v\n", err)
	}

	encryptedKubeconfig, err := ca.crypto.Encrypt(c.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to encrypt kubeconfig: %v\n", err)
	}
	err = ca.clusterRepository.UpdateClusterKubeconfig(ctx, c.Id, &encryptedKubeconfig)
	if err != nil {
		fmt.Printf("Failed to persist kubeconfig: %v\n", err)
	}

	err = ca.clusterRepository.UpdateClusterStatus(ctx, c.Id, entity.ClusterRunning)
	if err != nil {
		fmt.Printf("Failed to update cluster status: %v\n", err)
	}
}

func (ca *ClusterApplication) HandleClusterDeleted(c *coreValue.ClusterDeleted) {
	ctx := context.Background()
	err := ca.clusterRepository.DeleteCluster(ctx, c.Id)
	if err != nil {
		fmt.Printf("failed to delete cluster from database: %v\n", err)
		return
	}
}

package application

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/port"
	interfaces "starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/domain/value"
	"strings"
	"time"
)

type ClusterApplication struct {
	organizationRepository interfaces.OrganizationRepository
	clusterRepository      interfaces.ClusterRepository
	organizationService    *service.OrganizationService
	ssh                    port.SSH
	install                port.Install
	provision              port.Provision
	crypto                 port.Crypto
	queue                  port.Queue
}

func NewClusterApplication(
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	organizationService *service.OrganizationService,
	ssh port.SSH,
	install port.Install,
	provision port.Provision,
	crypto port.Crypto,
	queue port.Queue,
) *ClusterApplication {
	return &ClusterApplication{
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		organizationService:    organizationService,
		ssh:                    ssh,
		install:                install,
		provision:              provision,
		crypto:                 crypto,
		queue:                  queue,
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

	err = ca.queue.PublishCreateCluster(cluster)
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

	err = ca.queue.PublishDeleteCluster(cluster)
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (ca *ClusterApplication) HandleCreateCluster(c *entity.Cluster) {
	ctx := context.Background()
	publicKey, privateKey, err := ca.crypto.GenerateKeyPair()
	if err != nil {
		fmt.Printf("failed to generate ed25519 keypair: %v\n", err)
	}

	pubKeyStr := base64.StdEncoding.EncodeToString(publicKey)
	privKeyStr := base64.StdEncoding.EncodeToString(privateKey)

	encryptedPrivKeyStr, err := ca.crypto.Encrypt(privKeyStr)
	if err != nil {
		fmt.Printf("failed to encrypt private key: %v\n", err)
	}

	err = ca.clusterRepository.UpdateClusterPublicPrivateKey(ctx, c.Id, &pubKeyStr, &encryptedPrivKeyStr)
	if err != nil {
		fmt.Printf("failed to persist cluster public private key: %v\n", err)
	}

	organization, err := ca.organizationRepository.GetOrganization(ctx, c.OrganizationId)
	if err != nil {
		fmt.Printf("failed to get organization: %v", err)
	}

	trimmed := strings.TrimSpace(c.Name)
	clusterSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(organization.Name), clusterSlug)
	provisioningId, ip, err := ca.provision.ProvisionServer(ctx, projectName, publicKey)

	// Persist provisioningId regardless of outcome, to enable cleanup
	if provisioningId != "" {
		err = ca.clusterRepository.UpdateClusterPulumiStackId(ctx, c.Id, &provisioningId)
		if err != nil {
			fmt.Printf("failed to persist pulumi stack id: %v\n", err)
		}
	}

	if err != nil {
		fmt.Printf("failed to provision server: %v\n", err)
		return
	}

	err = ca.clusterRepository.UpdateClusterIPv4Address(ctx, c.Id, &ip)
	if err != nil {
		fmt.Printf("Failed to persist cluster ip address: %v\n", err)
	}

	err = ca.ssh.WaitForSSH(ip, 30*time.Second)
	if err != nil {
		fmt.Printf("SSH not available: %v\n", err)
		return
	}

	kubeconfig, err := ca.install.InstallK3s(ip, privateKey)
	kubeconfigBase64 := base64.StdEncoding.EncodeToString([]byte(kubeconfig))
	encryptedKubeconfig, err := ca.crypto.Encrypt(kubeconfigBase64)
	if err != nil {
		fmt.Printf("failed to encrypt kubeconfig: %v\n", err)
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

func (ca *ClusterApplication) HandleDeleteCluster(c *entity.Cluster) {
	ctx := context.Background()
	cluster, err := ca.clusterRepository.GetCluster(ctx, c.Id)
	if err != nil {
		fmt.Printf("failed to get cluster from database: %v\n", err)
		return
	}

	err = ca.provision.DeleteServer(ctx, *cluster.PulumiStackId)
	if err != nil {
		fmt.Printf("failed to delete server: %v\n", err)
		return
	}

	err = ca.clusterRepository.DeleteCluster(ctx, c.Id)
	if err != nil {
		fmt.Printf("failed to delete cluster from database: %v\n", err)
		return
	}
}

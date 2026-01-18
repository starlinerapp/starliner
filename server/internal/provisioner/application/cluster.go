package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"starliner.app/internal/api/domain/entity"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/repository/interface"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/provisioner/domain/port"
	"strings"
	"time"
)

type ClusterApplication struct {
	clusterRepository _interface.ClusterRepository
	ssh               port.SSH
	install           port.Install
	provision         port.Provision
	crypto            corePort.Crypto
}

func NewClusterApplication(
	clusterRepository _interface.ClusterRepository,
	ssh port.SSH,
	install port.Install,
	provision port.Provision,
	crypto corePort.Crypto,
) *ClusterApplication {
	return &ClusterApplication{
		clusterRepository: clusterRepository,
		ssh:               ssh,
		install:           install,
		provision:         provision,
		crypto:            crypto,
	}
}

func (ca *ClusterApplication) HandleProvisionCluster(c *value.Cluster) {
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

	if err != nil {
		fmt.Printf("failed to get organization: %v", err)
	}

	trimmed := strings.TrimSpace(c.Name)
	clusterSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(c.OrganizationName), clusterSlug)
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
	if err != nil {
		fmt.Printf("Failed to install k3s: %v\n", err)
		return
	}

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

func (ca *ClusterApplication) HandleDeleteCluster(c *value.Cluster) {
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

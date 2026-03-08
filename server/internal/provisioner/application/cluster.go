package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/service"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/provisioner/domain/port"
	"strings"
	"time"
)

type ClusterApplication struct {
	ssh               port.SSH
	install           port.Install
	provision         port.Provision
	queue             port.Queue
	crypto            corePort.Crypto
	normalizerService *service.NormalizerService
}

func NewClusterApplication(
	ssh port.SSH,
	install port.Install,
	provision port.Provision,
	queue port.Queue,
	crypto corePort.Crypto,
	normalizerService *service.NormalizerService,
) *ClusterApplication {
	return &ClusterApplication{
		ssh:               ssh,
		install:           install,
		provision:         provision,
		queue:             queue,
		crypto:            crypto,
		normalizerService: normalizerService,
	}
}

func (ca *ClusterApplication) HandleProvisionCluster(c *value.ProvisionCluster) {
	ctx := context.Background()
	publicKey, privateKey, err := ca.crypto.GenerateKeyPair()
	if err != nil {
		log.Printf("failed to generate ed25519 keypair: %v\n", err)
	}

	pubKeyStr := base64.StdEncoding.EncodeToString(publicKey)
	privKeyStr := base64.StdEncoding.EncodeToString(privateKey)

	clusterSlug, err := ca.normalizerService.FormatToDNS1123(c.Name)
	if err != nil {
		log.Printf("failed to normalize cluster name: %v\n", err)
	}

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(c.OrganizationName), clusterSlug)
	provisioningId, ip, err := ca.provision.ProvisionServer(ctx, c.ProvisioningCredential, projectName, publicKey)

	if err != nil {
		log.Printf("failed to provision server: %v\n", err)
		ca.HandleDeleteCluster(&value.DeleteCluster{
			Id:                     c.Id,
			ProvisioningId:         provisioningId,
			ProvisioningCredential: c.ProvisioningCredential,
		})
		return
	}

	err = ca.ssh.WaitForSSH(ip, 30*time.Second)
	if err != nil {
		log.Printf("SSH not available: %v\n", err)
		return
	}

	kubeconfig, err := ca.install.InstallK3s(ip, privateKey)
	if err != nil {
		log.Printf("Failed to install k3s: %v\n", err)
		return
	}

	kubeconfigBase64 := base64.StdEncoding.EncodeToString([]byte(kubeconfig))

	err = ca.queue.PublishClusterCreated(&value.ClusterCreated{
		Id:               c.Id,
		ProvisioningId:   provisioningId,
		IPv4Address:      ip,
		PublicKey:        pubKeyStr,
		PrivateKey:       privKeyStr,
		KubeconfigBase64: kubeconfigBase64,
	})

	if err != nil {
		log.Printf("failed to publish event: %v\n", err)
	}
}

func (ca *ClusterApplication) HandleDeleteCluster(c *value.DeleteCluster) {
	ctx := context.Background()
	err := ca.provision.DeleteServer(ctx, c.ProvisioningCredential, c.ProvisioningId)
	if err != nil {
		log.Printf("failed to delete server: %v\n", err)
		return
	}
	err = ca.queue.PublishClusterDeleted(&value.ClusterDeleted{
		Id: c.Id,
	})
	if err != nil {
		log.Printf("failed to publish cluster deleted event: %v\n", err)
	}
}

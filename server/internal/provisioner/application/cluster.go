package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/service"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/provisioner/domain/port"
)

type ClusterApplication struct {
	ssh               port.SSH
	install           port.Install
	provision         port.Provision
	queue             port.Queue
	crypto            corePort.Crypto
	logPublisher      port.LogPublisher
	normalizerService *service.NormalizerService
}

func NewClusterApplication(
	ssh port.SSH,
	install port.Install,
	provision port.Provision,
	queue port.Queue,
	crypto corePort.Crypto,
	logPublisher port.LogPublisher,
	normalizerService *service.NormalizerService,
) *ClusterApplication {
	return &ClusterApplication{
		ssh:               ssh,
		install:           install,
		provision:         provision,
		queue:             queue,
		crypto:            crypto,
		logPublisher:      logPublisher,
		normalizerService: normalizerService,
	}
}

func (ca *ClusterApplication) HandleProvisionCluster(c *value.ProvisionCluster) {
	defer func() {
		if ca.logPublisher == nil {
			return
		}
		if err := ca.logPublisher.PublishLogEnd(c.Id); err != nil {
			log.Printf("failed to publish log end: %v", err)
		}
	}()

	var logBuf strings.Builder
	appendStatus := func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		logBuf.WriteString(line)
		if ca.logPublisher == nil {
			return
		}
		if err := ca.logPublisher.PublishLogChunk(c.Id, []byte(line)); err != nil {
			log.Printf("failed to publish log chunk: %v", err)
		}
	}

	ctx := context.Background()

	publicKey, privateKey, err := ca.crypto.GenerateKeyPair()
	if err != nil {
		appendStatus("==> ERROR: failed to generate ed25519 keypair: %v\n", err)
		log.Printf("failed to generate ed25519 keypair: %v\n", err)
		return
	}

	pubKeyStr := base64.StdEncoding.EncodeToString(publicKey)
	privKeyStr := base64.StdEncoding.EncodeToString(privateKey)

	clusterSlug, err := ca.normalizerService.FormatToDNS1123(c.Name)
	if err != nil {
		appendStatus("==> ERROR: failed to normalize cluster name: %v\n", err)
		log.Printf("failed to normalize cluster name: %v\n", err)
		return
	}

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(c.OrganizationName), clusterSlug)

	appendStatus("==> Provisioning server %q...\n", projectName)
	provisioningId, ip, provisionLogs, err := ca.provision.ProvisionServer(ctx, c.Id, c.ProvisioningCredential, projectName, c.ServerType, publicKey)
	logBuf.WriteString(provisionLogs)
	if err != nil {
		appendStatus("==> ERROR: failed to provision server: %v\n", err)
		log.Printf("failed to provision server: %v\n", err)
		ca.HandleDeleteCluster(&value.DeleteCluster{
			Id:                     c.Id,
			ProvisioningId:         provisioningId,
			ProvisioningCredential: c.ProvisioningCredential,
		})
		return
	}
	appendStatus("==> Server provisioned at %s\n", ip)

	appendStatus("==> Waiting for SSH...\n")

	pemBytes, err := ca.crypto.EncodePrivateKeyToPEM(privateKey)
	if err != nil {
		appendStatus("==> ERROR: failed to encode private key to PEM: %v\n", err)
		log.Printf("failed to encode private key to PEM: %v\n", err)
		ca.HandleDeleteCluster(&value.DeleteCluster{
			Id:                     c.Id,
			ProvisioningId:         provisioningId,
			ProvisioningCredential: c.ProvisioningCredential,
		})
		return
	}

	if err := ca.ssh.WaitForSSH(ip, "root", pemBytes, 30*time.Second); err != nil {
		appendStatus("==> ERROR: SSH not available: %v\n", err)
		log.Printf("SSH not available: %v\n", err)
		return
	}
	appendStatus("==> SSH is ready\n")

	appendStatus("==> Installing K3s...\n")
	kubeconfig, installLogs, err := ca.install.InstallK3s(c.Id, ip, privateKey)
	logBuf.WriteString(installLogs)
	if err != nil {
		appendStatus("==> ERROR: failed to install k3s: %v\n", err)
		log.Printf("Failed to install k3s: %v\n", err)
		return
	}
	appendStatus("==> K3s installed\n")

	kubeconfigBase64 := base64.StdEncoding.EncodeToString([]byte(kubeconfig))

	err = ca.queue.PublishClusterCreated(&value.ClusterCreated{
		Id:               c.Id,
		ProvisioningId:   provisioningId,
		IPv4Address:      ip,
		PublicKey:        pubKeyStr,
		PrivateKey:       privKeyStr,
		KubeconfigBase64: kubeconfigBase64,
		Logs:             logBuf.String(),
	})

	if err != nil {
		appendStatus("==> ERROR: failed to publish cluster created event: %v\n", err)
		log.Printf("failed to publish event: %v\n", err)
	}
}

func (ca *ClusterApplication) HandleDeleteCluster(c *value.DeleteCluster) {
	defer func() {
		if ca.logPublisher == nil {
			return
		}
		if err := ca.logPublisher.PublishLogEnd(c.Id); err != nil {
			log.Printf("failed to publish log end: %v", err)
		}
	}()

	appendStatus := func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		if ca.logPublisher == nil {
			return
		}
		if err := ca.logPublisher.PublishLogChunk(c.Id, []byte(line)); err != nil {
			log.Printf("failed to publish log chunk: %v", err)
		}
	}

	ctx := context.Background()

	appendStatus("==> Deleting server...\n")
	if err := ca.provision.DeleteServer(ctx, c.Id, c.ProvisioningCredential, c.ProvisioningId); err != nil {
		appendStatus("==> ERROR: failed to delete server: %v\n", err)
		log.Printf("failed to delete server: %v\n", err)
		return
	}
	appendStatus("==> Server deleted\n")

	if err := ca.queue.PublishClusterDeleted(&value.ClusterDeleted{
		Id: c.Id,
	}); err != nil {
		appendStatus("==> ERROR: failed to publish cluster deleted event: %v\n", err)
		log.Printf("failed to publish cluster deleted event: %v\n", err)
	}
}

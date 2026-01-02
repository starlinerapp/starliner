package cluster

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/google/uuid"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"go.uber.org/fx"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"os/exec"
	"starliner.app/internal/config"
	"starliner.app/internal/core/cluster/ansible"
	"starliner.app/internal/crypto"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/infrastructure/queue/proto/v1"
	interfaces "starliner.app/internal/repository/interface"
	"strings"
	"time"
)

type Orchestrator struct {
	cfg                    *config.Config
	organizationRepository interfaces.OrganizationRepository
	clusterRepository      interfaces.ClusterRepository
	clusterSubscriber      *queue.Subscriber[*v1.Cluster]
}

func RegisterOrchestrator(lc fx.Lifecycle, o *Orchestrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewOrchestrator(
	cfg *config.Config,
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	clusterSubscriber *queue.Subscriber[*v1.Cluster],
) *Orchestrator {
	return &Orchestrator{
		cfg:                    cfg,
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		clusterSubscriber:      clusterSubscriber,
	}
}

func (o *Orchestrator) Start() error {
	go func() {
		err := o.clusterSubscriber.Subscribe(queue.CreateCluster, "*", "createCluster", o.handleCreateCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	go func() {
		err := o.clusterSubscriber.Subscribe(queue.DeleteCluster, "*", "deleteCluster", o.handleDeleteCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}

func (o *Orchestrator) handleCreateCluster(c *v1.Cluster) {
	ctx := context.Background()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("failed to generate ed25519 keypair: %v\n", err)
	}

	pubKeyStr := base64.StdEncoding.EncodeToString(publicKey)
	privKeyStr := base64.StdEncoding.EncodeToString(privateKey)

	encryptionKey, err := base64.StdEncoding.DecodeString(o.cfg.EncryptionKeyBase64)
	if err != nil {
		fmt.Printf("failed to decode encryption key: %v\n", err)
	}
	encryptedPrivKeyStr, err := crypto.Encrypt(privKeyStr, encryptionKey)
	if err != nil {
		fmt.Printf("failed to encrypt private key: %v\n", err)
	}

	err = o.clusterRepository.UpdateClusterPublicPrivateKey(ctx, c.Id, &pubKeyStr, &encryptedPrivKeyStr)
	if err != nil {
		fmt.Printf("failed to persist cluster public private key: %v\n", err)
	}

	organization, err := o.organizationRepository.GetOrganization(ctx, c.OrganizationId)
	if err != nil {
		fmt.Printf("failed to get organization: %v", err)
	}

	trimmed := strings.TrimSpace(c.Name)
	clusterSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(organization.Name), clusterSlug)
	stackName := auto.FullyQualifiedStackName("organization", projectName, uuid.New().String())

	err = o.clusterRepository.UpdateClusterPulumiStackId(ctx, c.Id, &stackName)
	if err != nil {
		fmt.Printf("failed to persist pulumi stack id: %v\n", err)
	}

	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployFunc(DeployParams{
		ServerName: projectName,
		publicKey:  publicKey,
	}))
	if err != nil {
		fmt.Printf("failed to set up a workspace: %v\n", err)
		return
	}

	w := s.Workspace()

	err = w.InstallPlugin(ctx, "hcloud", "1.29")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		return
	}

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		return
	}

	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		return
	}

	ip, ok := res.Outputs["serverIp"].Value.(string)
	if !ok {
		fmt.Println("Failed to unmarshall output")
		return
	}

	err = o.clusterRepository.UpdateClusterIPv4Address(ctx, c.Id, &ip)
	if err != nil {
		fmt.Printf("Failed to persist cluster ip address: %v\n", err)
	}

	err = waitForSSH(ip, 30*time.Second)
	if err != nil {
		fmt.Printf("SSH isn't available: %v\n", err)
		return
	}

	// Install k3s
	tmpPlaybook, err := os.CreateTemp("", "ansible-*.yml")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPlaybook.Name())

	_, err = tmpPlaybook.Write([]byte(ansible.K3sPlaybook))
	if err != nil {
		fmt.Printf("Failed to write to temp file: %v\n", err)
	}
	err = tmpPlaybook.Close()
	if err != nil {
		return
	}

	tmpPrivateKey, err := os.CreateTemp("", "private-key-*.pem")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPrivateKey.Name())

	block, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		fmt.Printf("failed to marshal private key: %v\n", err)
		return
	}

	pemBytes := pem.EncodeToMemory(block)
	if pemBytes == nil {
		fmt.Println("failed to encode private key to PEM")
		return
	}

	_, err = tmpPrivateKey.Write(pemBytes)
	if err != nil {
		fmt.Printf("failed to write private key: %v\n", err)
		return
	}
	err = tmpPrivateKey.Close()
	if err != nil {
		return
	}

	args := []string{
		tmpPlaybook.Name(),
		"-i", fmt.Sprintf("%s,", ip),
		"-u", "root",
		"--private-key", tmpPrivateKey.Name(),
		"-e", "ansible_ssh_common_args='-o StrictHostKeyChecking=no'",
		"-e", fmt.Sprintf("target_host=%s", ip),
	}

	out, err := exec.Command("ansible-playbook", args...).CombinedOutput()
	fmt.Printf("Ansible output:\n%s\n", string(out))

	if err != nil {
		fmt.Printf("Failed to install k3s: %v\n", err)
		return
	}

	err = o.clusterRepository.UpdateClusterStatus(ctx, c.Id, domain.ClusterRunning)
	if err != nil {
		fmt.Printf("Failed to update cluster status: %v\n", err)
	}
}

func (o *Orchestrator) handleDeleteCluster(c *v1.Cluster) {
	ctx := context.Background()

	cluster, err := o.clusterRepository.GetCluster(ctx, c.Id)
	if err != nil {
		fmt.Printf("failed to get cluster from database: %v\n", err)
		return
	}

	parts := strings.Split(*cluster.PulumiStackId, "/")
	projectName := parts[1]

	pubKeyBytes, err := base64.StdEncoding.DecodeString(*cluster.PublicKey)
	if err != nil {
		fmt.Printf("failed to decode public key: %v\n", err)
		return
	}

	s, err := auto.SelectStackInlineSource(ctx, *cluster.PulumiStackId, projectName, deployFunc(DeployParams{
		ServerName: projectName,
		publicKey:  pubKeyBytes,
	}))
	if err != nil {
		fmt.Printf("failed to select stack for deletion: %v\n", err)
		return
	}

	w := s.Workspace()
	err = w.InstallPlugin(ctx, "hcloud", "1.29")
	if err != nil {
		fmt.Printf("failed to install program plugins: %v\n", err)
	}

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("failed to refresh stack: %v\n", err)
	}

	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)
	_, err = s.Destroy(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("failed to destroy stack: %v\n", err)
	}

	err = o.clusterRepository.DeleteCluster(ctx, c.Id)
	if err != nil {
		fmt.Printf("failed to delete cluster from database: %v\n", err)
		return
	}
}

func waitForSSH(ip string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "22"), 5*time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for ssh on %s", ip)
		}
		time.Sleep(5 * time.Second)
	}
}

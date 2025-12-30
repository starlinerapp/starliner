package cluster

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"go.uber.org/fx"
	"log"
	"os"
	"starliner.app/pkg/config"
	"starliner.app/pkg/crypto"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
	interfaces "starliner.app/pkg/repository/interface"
	"strings"
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
	err = o.clusterRepository.UpdateClusterPublicPrivateKey(ctx, c.Id, &pubKeyStr, &encryptedPrivKeyStr)
	if err != nil {
		fmt.Printf("failed to persist cluster public private key: %v\n", err)
	}

	organization, err := o.organizationRepository.GetOrganization(ctx, c.OrganizationId)
	if err != nil {
		fmt.Printf("failed to get organization: %v", err)
	}

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(organization.Name), c.Name)
	stackName := auto.FullyQualifiedStackName("organization", projectName, uuid.New().String())

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
}

func (o *Orchestrator) handleDeleteCluster(c *v1.Cluster) {
	ctx := context.Background()
	err := o.clusterRepository.DeleteCluster(ctx, c.Id)
	if err != nil {
		fmt.Printf("failed to delete cluster from database: %v", err)
	}
}

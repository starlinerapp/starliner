package cluster

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"go.uber.org/fx"
	"log"
	"os"
	"starliner.app/pkg/config"
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
	fmt.Printf("Create cluster %s\n", c.Name)
	organization, err := o.organizationRepository.GetOrganization(ctx, c.OrganizationId)
	if err != nil {
		fmt.Printf("failed to get organization: %v", err)
	}

	projectName := fmt.Sprintf("%s-%s", strings.ToLower(organization.Name), c.Name)
	stackName := auto.FullyQualifiedStackName("organization", projectName, uuid.New().String())

	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployFunc)
	if err != nil {
		fmt.Printf("Failed to set up a workspace: %v\n", err)
		return
	}
	fmt.Printf("Created/Selected stack %q\n", stackName)

	w := s.Workspace()

	fmt.Println("Installing Hetzner Cloud plugin")
	err = w.InstallPlugin(ctx, "hcloud", "1.29")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		return
	}
	fmt.Println("Successfully installed Hetzner Cloud plugin")

	fmt.Println("Starting refresh")
	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		return
	}
	fmt.Println("Refresh succeeded!")

	fmt.Println("Starting update")
	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		return
	}
	fmt.Println("Update succeeded!")

	ip, ok := res.Outputs["serverIp"].Value.(string)
	if !ok {
		fmt.Println("Failed to unmarshall output")
		return
	}

	fmt.Printf("Server ipv4: %s\n", ip)
}

func (o *Orchestrator) handleDeleteCluster(c *v1.Cluster) {
	ctx := context.Background()
	if c.Id == nil {
		fmt.Printf("cannot delete cluster: missing cluster id")
		return
	}
	err := o.clusterRepository.DeleteCluster(ctx, *c.Id)
	if err != nil {
		fmt.Printf("failed to delete cluster from database: %v", err)
	}
}

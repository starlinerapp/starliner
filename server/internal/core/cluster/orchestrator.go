package cluster

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"go.uber.org/fx"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"io/fs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/core/cluster/helm"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/infrastructure/queue/proto/v1"
	interfaces "starliner.app/internal/repository/interface"
	"starliner.app/internal/service"
)

type Orchestrator struct {
	clusterRepository interfaces.ClusterRepository
	cryptoService     *service.CryptoService
	projectSubscriber *queue.Subscriber[*v1.Project]
}

func RegisterOrchestrator(lc fx.Lifecycle, o *Orchestrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewOrchestrator(
	clusterRepository interfaces.ClusterRepository,
	cryptoService *service.CryptoService,
	projectSubscriber *queue.Subscriber[*v1.Project],
) *Orchestrator {
	return &Orchestrator{
		clusterRepository: clusterRepository,
		cryptoService:     cryptoService,
		projectSubscriber: projectSubscriber,
	}
}

func (o *Orchestrator) Start() error {
	go func() {
		err := o.projectSubscriber.Subscribe(queue.CreateProject, "*", "createProject", o.handleCreateProject)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}

func (o *Orchestrator) handleCreateProject(p *v1.Project) {
	ctx := context.Background()

	cluster, err := o.clusterRepository.GetCluster(ctx, p.ClusterId)
	if err != nil {
		fmt.Printf("failed to get cluster from database: %v\n", err)
		return
	}

	kubeconfigBase64, err := o.cryptoService.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		fmt.Printf("failed to decrypt kubeconfig: %v\n", err)
		return
	}
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfigBase64)
	if err != nil {
		fmt.Printf("failed to decode kubeconfig: %v\n", err)
		return
	}

	tmpDir, err := os.MkdirTemp("", "helm-chart-*")
	if err != nil {
		fmt.Printf("failed to create temp directory: %v\n", err)
		return
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("failed to remove temp directory: %v\n", err)
		}
	}()

	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, kubeconfigBytes, 0600)
	if err != nil {
		fmt.Printf("failed to write kubeconfig: %v\n", err)
		return
	}

	err = fs.WalkDir(helm.Chart, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("template", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(tmpDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := helm.Chart.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0644)
	})
	if err != nil {
		fmt.Printf("failed to copy helm chart: %v\n", err)
		return
	}

	chart, err := loader.Load(tmpDir)
	if err != nil {
		fmt.Printf("failed to load helm chart: %v\n", err)
		return
	}

	configFlags := genericclioptions.NewConfigFlags(false)
	configFlags.KubeConfig = &kubeconfigPath

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		configFlags,
		"default",
		"secret",
	); err != nil {
		fmt.Printf("failed to initialize helm action configuration: %v\n", err)
		return
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = "test-release"
	install.Namespace = "default"
	install.WaitStrategy = "watcher"

	fmt.Printf("Installing helm chart")
	_, err = install.Run(chart, map[string]interface{}{
		"ingress": map[string]interface{}{
			"host": fmt.Sprintf("%s.nip.io", *cluster.IPv4Address),
		},
	})
	if err != nil {
		fmt.Printf("failed to install helm chart: %v\n", err)
		return
	}
}

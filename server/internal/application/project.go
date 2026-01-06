package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	"io/fs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/domain/entity"
	interfaces "starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/domain/value"
	"starliner.app/internal/infrastructure/crypto"
	"starliner.app/internal/infrastructure/helm"
	"starliner.app/internal/infrastructure/queue"
	v1 "starliner.app/internal/infrastructure/queue/proto/v1"
	"strconv"
	"strings"
)

type ProjectApplication struct {
	crypto                 *crypto.Crypto
	organizationService    *service.OrganizationService
	projectRepository      interfaces.ProjectRepository
	clusterRepository      interfaces.ClusterRepository
	organizationRepository interfaces.OrganizationRepository
	environmentRepository  interfaces.EnvironmentRepository
	projectPublisher       *queue.Publisher[*v1.Project]
}

func NewProjectApplication(
	crypto *crypto.Crypto,
	organizationService *service.OrganizationService,
	projectRepository interfaces.ProjectRepository,
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	environmentRepository interfaces.EnvironmentRepository,
	projectPublisher *queue.Publisher[*v1.Project],
) *ProjectApplication {
	return &ProjectApplication{
		crypto:                 crypto,
		organizationService:    organizationService,
		projectRepository:      projectRepository,
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		environmentRepository:  environmentRepository,
		projectPublisher:       projectPublisher,
	}
}

func (ps *ProjectApplication) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64) (*value.Project, error) {
	err := ps.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	project, err := ps.projectRepository.CreateProject(ctx, name, organizationId, clusterId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	environment, err := ps.environmentRepository.CreateEnvironment(ctx, productionEnvName, strings.ToLower(productionEnvName), project.Id)
	if err != nil {
		return nil, err
	}

	err = ps.projectPublisher.Publish(queue.CreateProject, strconv.FormatInt(project.Id, 10), &v1.Project{
		Id:             project.Id,
		Name:           project.Name,
		OrganizationId: project.OrganizationId,
		ClusterId:      *project.ClusterId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	environmentsModel := value.NewEnvironments([]*entity.Environment{environment})
	projectModel := value.NewProject(project)
	projectModel.Environments = environmentsModel

	return projectModel, nil
}

func (ps *ProjectApplication) GetProject(ctx context.Context, projectId int64, userId int64) (*value.Project, error) {
	project, err := ps.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProject(project), nil
}

func (ps *ProjectApplication) DeleteProject(ctx context.Context, projectId int64, userId int64) error {
	return ps.projectRepository.DeleteProject(ctx, projectId, userId)
}

func (ps *ProjectApplication) HandleCreateProject(p *v1.Project) {
	ctx := context.Background()

	cluster, err := ps.clusterRepository.GetCluster(ctx, p.ClusterId)
	if err != nil {
		fmt.Printf("failed to get cluster from database: %v\n", err)
		return
	}

	kubeconfigBase64, err := ps.crypto.Decrypt(*cluster.Kubeconfig)
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

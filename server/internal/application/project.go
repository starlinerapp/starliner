package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/port"
	interfaces "starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/domain/value"
	"strings"
)

type ProjectApplication struct {
	organizationService    *service.OrganizationService
	projectService         *service.ProjectService
	projectRepository      interfaces.ProjectRepository
	clusterRepository      interfaces.ClusterRepository
	organizationRepository interfaces.OrganizationRepository
	environmentRepository  interfaces.EnvironmentRepository
	deploy                 port.Deploy
	crypto                 port.Crypto
	queue                  port.Queue
}

func NewProjectApplication(
	organizationService *service.OrganizationService,
	projectService *service.ProjectService,
	projectRepository interfaces.ProjectRepository,
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	environmentRepository interfaces.EnvironmentRepository,
	deploy port.Deploy,
	crypto port.Crypto,
	queue port.Queue,
) *ProjectApplication {
	return &ProjectApplication{
		organizationService:    organizationService,
		projectService:         projectService,
		projectRepository:      projectRepository,
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		environmentRepository:  environmentRepository,
		deploy:                 deploy,
		crypto:                 crypto,
		queue:                  queue,
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

	err = ps.queue.PublishCreateProject(project)
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	environmentsModel := value.NewEnvironments([]*entity.Environment{environment})
	projectModel := value.NewProject(project)
	projectModel.Environments = environmentsModel

	return projectModel, nil
}

func (ps *ProjectApplication) UpdateProjectName(ctx context.Context, projectName string, projectId int64, userId int64) error {
	if err := ps.projectService.ValidateUserHasPermission(ctx, projectId, userId); err != nil {
		return err
	}
	return ps.projectRepository.UpdateProjectName(ctx, projectName, projectId)
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

func (ps *ProjectApplication) HandleCreateProject(p *entity.ProjectWithEnvironments) {
	ctx := context.Background()

	cluster, err := ps.clusterRepository.GetCluster(ctx, *p.ClusterId)
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

	tmpDir, err := os.MkdirTemp("", "kubeconfig-*")
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

	err = ps.deploy.DeployNginx(*cluster.IPv4Address, kubeconfigPath)
	if err != nil {
		fmt.Printf("failed to install helm chart: %v\n", err)
		return
	}
}

package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
)

type DeploymentRepository interface {
	CreateGitDeployment(
		ctx context.Context,
		environmentId int64,
		serviceName string,
		port string,
		gitUrl string,
		projectRepositoryPath string,
		dockerfilePath string,
		envs []*value.EnvVar,
	) (deployment *entity.GitDeployment, err error)

	UpdateGitDeployment(
		ctx context.Context,
		deploymentId int64,
		port string,
		projectRepositoryPath string,
		dockerfilePath string,
		envs []*value.EnvVar,
	) (deployment *entity.GitDeployment, err error)

	CreateImageDeployment(
		ctx context.Context,
		serviceName string,
		imageName string,
		tag string,
		port string,
		environmentId int64,
		envs []*value.EnvVar,
	) (deployment *entity.ImageDeployment, err error)

	UpdateImageDeployment(
		ctx context.Context,
		deploymentId int64,
		imageName string,
		tag string,
		port string,
		envs []*value.EnvVar,
	) (deployment *entity.ImageDeployment, err error)

	CreateIngressDeployment(
		ctx context.Context,
		serviceName string,
		port string,
		environmentId int64,
		hosts []*value.IngressHost,
	) (*entity.IngressDeployment, error)

	UpdateIngressDeployment(
		ctx context.Context,
		deploymentId int64,
		port string,
		environmentId int64,
		hosts []*value.IngressHost,
	) (*entity.IngressDeployment, error)

	CreateDatabaseDeployment(
		ctx context.Context,
		name string,
		port string,
		environmentId int64,
	) (deployment *entity.DatabaseDeployment, err error)

	UpdateDatabaseDeploymentCredentials(
		ctx context.Context,
		dbName string,
		deploymentId int64,
		username string,
		password string,
	) error

	GetDeploymentEnvs(ctx context.Context, deploymentId int64) ([]*entity.EnvVar, error)

	GetUserDeployment(ctx context.Context, userId int64, deploymentId int64) (*entity.Deployment, error)

	GetDeploymentWithNamespace(ctx context.Context, deploymentId int64) (*entity.Deployment, error)

	GetDeploymentCluster(ctx context.Context, deploymentId int64) (*entity.Cluster, error)

	DeleteDeployment(ctx context.Context, deploymentId int64) error

	GetAllDeploymentsWithKubeconfig(ctx context.Context) ([]*entity.DeploymentWithKubeconfig, error)

	UpdateDeploymentStatus(ctx context.Context, deploymentId int64, status string) error

	GetEnvironmentDeploymentByName(ctx context.Context, environmentId int64, serviceName string) (*entity.Deployment, error)
}

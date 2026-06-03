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
		args []*value.Arg,
	) (deployment *entity.GitDeployment, err error)

	UpdateGitDeployment(
		ctx context.Context,
		deploymentId int64,
		port string,
		projectRepositoryPath string,
		dockerfilePath string,
		envs []*value.EnvVar,
		args []*value.Arg,
	) (deployment *entity.GitDeployment, err error)

	CreateImageDeployment(
		ctx context.Context,
		serviceName string,
		imageName string,
		tag string,
		port string,
		volumeSizeMiB *int32,
		volumeMountPath *string,
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

	GetGitDeploymentArgs(ctx context.Context, deploymentId int64) ([]*entity.Arg, error)

	GetUserDeployment(ctx context.Context, userId int64, deploymentId int64) (*entity.Deployment, error)

	GetDeploymentWithNamespace(ctx context.Context, deploymentId int64) (*entity.Deployment, error)

	GetDeploymentCluster(ctx context.Context, deploymentId int64) (*entity.Cluster, error)

	SoftDeleteDeploymentVolume(ctx context.Context, deploymentId int64) error

	SoftDeleteDeployment(ctx context.Context, deploymentId int64) error

	RepointIngressPathsTargetDeployment(ctx context.Context, oldDeploymentId int64, newDeploymentId int64) error

	SoftDeleteDeploymentsByEnvironmentId(ctx context.Context, environmentId int64) error

	GetAllDeploymentsWithKubeconfig(ctx context.Context) ([]*entity.DeploymentWithKubeconfig, error)

	UpdateDeploymentStatus(ctx context.Context, deploymentId int64, status string) error

	GetDeploymentStatusLogs(ctx context.Context, userId int64, deploymentId int64) (*entity.DeploymentStatusLogs, error)
	SetDeploymentStatusLogs(ctx context.Context, deploymentId int64, logs string, rolloutStatus string) error

	GetEnvironmentDeploymentByName(ctx context.Context, environmentId int64, serviceName string) (*entity.Deployment, error)

	GetIngressHostByName(ctx context.Context, hostName string) (*entity.IngressHostDeployment, error)

	IsIngressDeployment(ctx context.Context, deploymentId int64) (bool, error)

	GetGitDeploymentsByRepositoryUrl(ctx context.Context, repositoryUrl string) ([]*entity.GitDeployment, error)

	GetUserGitDeploymentById(ctx context.Context, userId int64, deploymentId int64) (*entity.GitDeployment, error)

	GetUserImageDeploymentById(ctx context.Context, userId int64, deploymentId int64) (*entity.ImageDeployment, error)

	GetUserDatabaseDeploymentById(ctx context.Context, userId int64, deploymentId int64) (*entity.DatabaseDeployment, error)
}

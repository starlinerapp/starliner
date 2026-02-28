package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
)

type DeploymentRepository interface {
	CreateImageDeployment(
		ctx context.Context,
		serviceName string,
		imageName string,
		tag string,
		port string,
		status string,
		environmentId int64,
	) (deployment *entity.ImageDeployment, err error)

	CreateIngressDeployment(
		ctx context.Context,
		serviceName string,
		port string,
		status string,
		environmentId int64,
	) (*entity.IngressDeployment, error)

	CreateDatabaseDeployment(
		ctx context.Context,
		name string,
		port string,
		status string,
		username string,
		password string,
		environmentId int64,
	) (deployment *entity.DatabaseDeployment, err error)

	GetUserDeployment(ctx context.Context, userId int64, deploymentId int64) (*entity.Deployment, error)

	GetDeploymentCluster(ctx context.Context, deploymentId int64) (*entity.Cluster, error)

	DeleteDeployment(ctx context.Context, deploymentId int64) error

	GetAllDeploymentsWithKubeconfig(ctx context.Context) ([]*entity.DeploymentWithKubeconfig, error)

	UpdateDeploymentStatus(ctx context.Context, deploymentId int64, status string) error
}

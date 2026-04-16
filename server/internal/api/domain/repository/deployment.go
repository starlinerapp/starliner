package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type DeploymentRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

var _ interfaces.DeploymentRepository = (*DeploymentRepository)(nil)

func NewDeploymentRepository(db *sql.DB, queries *sqlc.Queries) interfaces.DeploymentRepository {
	return &DeploymentRepository{db: db, queries: queries}
}

func (dr *DeploymentRepository) CreateGitDeployment(
	ctx context.Context,
	environmentId int64,
	serviceName string,
	port string,
	gitUrl string,
	projectRepositoryPath string,
	dockerfilePath string,
	envs []*value.EnvVar,
	args []*value.Arg,
) (deployment *entity.GitDeployment, err error) {
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := dr.queries.WithTx(tx)

	d, err := qtx.CreateGitDeployment(ctx, sqlc.CreateGitDeploymentParams{
		Name:           serviceName,
		Port:           port,
		EnvironmentID:  environmentId,
		Url:            gitUrl,
		ProjectPath:    projectRepositoryPath,
		DockerfilePath: dockerfilePath,
	})
	if err != nil {
		return nil, err
	}

	vars := make([]*entity.EnvVar, len(envs))
	for i, e := range envs {
		variable, err := qtx.CreateDeploymentEnvVar(ctx, sqlc.CreateDeploymentEnvVarParams{
			DeploymentID: d.DeploymentID,
			Name:         e.Name,
			Value:        e.Value,
		})
		if err != nil {
			return nil, err
		}

		vars[i] = &entity.EnvVar{
			Name:  variable.Name,
			Value: variable.Value,
		}
	}

	resultArgs := make([]*entity.Arg, len(args))
	for i, a := range args {
		arg, err := qtx.CreateDeploymentArg(ctx, sqlc.CreateDeploymentArgParams{
			DeploymentID: d.DeploymentID,
			Name:         a.Name,
			Value:        a.Value,
		})
		if err != nil {
			return nil, err
		}

		resultArgs[i] = &entity.Arg{
			Name:  arg.Name,
			Value: arg.Value,
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.GitDeployment{
		Id:                    d.DeploymentID,
		Name:                  d.Name,
		Port:                  d.Port,
		EnvironmentId:         d.EnvironmentID,
		GitUrl:                d.Url,
		ProjectRepositoryPath: d.ProjectPath,
		DockerfilePath:        d.DockerfilePath,
		EnvVars:               vars,
		Args:                  resultArgs,
	}, nil
}

func (dr *DeploymentRepository) UpdateGitDeployment(
	ctx context.Context,
	deploymentId int64,
	port string,
	projectRepositoryPath string,
	dockerfilePath string,
	envs []*value.EnvVar,
	args []*value.Arg,
) (deployment *entity.GitDeployment, err error) {
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := dr.queries.WithTx(tx)

	d, err := qtx.UpdateGitDeployment(ctx, sqlc.UpdateGitDeploymentParams{
		Port:           port,
		DeploymentID:   deploymentId,
		ProjectPath:    projectRepositoryPath,
		DockerfilePath: dockerfilePath,
	})
	if err != nil {
		return nil, err
	}

	if err := qtx.DeleteEnvVarsByDeploymentId(ctx, deploymentId); err != nil {
		return nil, err
	}

	vars := make([]*entity.EnvVar, len(envs))
	for i, e := range envs {
		variable, err := qtx.CreateDeploymentEnvVar(ctx, sqlc.CreateDeploymentEnvVarParams{
			DeploymentID: d.DeploymentID,
			Name:         e.Name,
			Value:        e.Value,
		})
		if err != nil {
			return nil, err
		}

		vars[i] = &entity.EnvVar{
			Name:  variable.Name,
			Value: variable.Value,
		}
	}

	if err := qtx.DeleteArgsByDeploymentId(ctx, deploymentId); err != nil {
		return nil, err
	}

	resultArgs := make([]*entity.Arg, len(args))
	for i, a := range args {
		arg, err := qtx.CreateDeploymentArg(ctx, sqlc.CreateDeploymentArgParams{
			DeploymentID: d.DeploymentID,
			Name:         a.Name,
			Value:        a.Value,
		})
		if err != nil {
			return nil, err
		}

		resultArgs[i] = &entity.Arg{
			Name:  arg.Name,
			Value: arg.Value,
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.GitDeployment{
		Id:                    d.DeploymentID,
		Name:                  d.ServiceName,
		Status:                string(d.Status),
		Port:                  d.Port,
		EnvironmentId:         d.EnvironmentID,
		GitUrl:                d.Url,
		ProjectRepositoryPath: d.ProjectPath,
		DockerfilePath:        d.DockerfilePath,
		EnvVars:               vars,
		Args:                  resultArgs,
	}, nil
}

func (dr *DeploymentRepository) CreateImageDeployment(
	ctx context.Context,
	serviceName string,
	imageName string,
	tag string,
	port string,
	volumeSizeMiB *int32,
	volumeMountPath *string,
	environmentId int64,
	envs []*value.EnvVar,
) (deployment *entity.ImageDeployment, err error) {
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := dr.queries.WithTx(tx)
	d, err := qtx.CreateImageDeployment(ctx, sqlc.CreateImageDeploymentParams{
		Port:          port,
		EnvironmentID: environmentId,
		Tag:           tag,
		ServiceName:   serviceName,
		ImageName:     imageName,
	})
	if err != nil {
		return nil, err
	}

	var resultVolumeSizeMiB *int32
	var resultVolumeMountPath *string
	if volumeSizeMiB != nil {
		mountPath := "/data"
		if volumeMountPath != nil && *volumeMountPath != "" {
			mountPath = *volumeMountPath
		}
		dv, err := qtx.CreateDeploymentVolume(ctx, sqlc.CreateDeploymentVolumeParams{
			DeploymentID:  sql.NullInt64{Int64: d.DeploymentID, Valid: true},
			VolumeSizeMib: *volumeSizeMiB,
			MountPath:     mountPath,
		})
		if err != nil {
			return nil, err
		}
		resultVolumeSizeMiB = &dv.VolumeSizeMib
		resultVolumeMountPath = &dv.MountPath
	}

	vars := make([]*entity.EnvVar, len(envs))
	for i, e := range envs {
		variable, err := qtx.CreateDeploymentEnvVar(ctx, sqlc.CreateDeploymentEnvVarParams{
			DeploymentID: d.DeploymentID,
			Name:         e.Name,
			Value:        e.Value,
		})
		if err != nil {
			return nil, err
		}

		vars[i] = &entity.EnvVar{
			Name:  variable.Name,
			Value: variable.Value,
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &entity.ImageDeployment{
		Id:              d.DeploymentID,
		Status:          string(d.Status),
		ServiceName:     d.ServiceName,
		ImageName:       d.ImageName,
		Tag:             d.ImageTag,
		Port:            d.Port,
		VolumeSizeMiB:   resultVolumeSizeMiB,
		VolumeMountPath: resultVolumeMountPath,
		EnvironmentId:   d.EnvironmentID,
		EnvVars:         vars,
	}, nil
}

func (dr *DeploymentRepository) UpdateImageDeployment(
	ctx context.Context,
	deploymentId int64,
	imageName string,
	tag string,
	port string,
	envs []*value.EnvVar,
) (deployment *entity.ImageDeployment, err error) {
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := dr.queries.WithTx(tx)

	d, err := qtx.UpdateImageDeployment(ctx, sqlc.UpdateImageDeploymentParams{
		Port:         port,
		DeploymentID: deploymentId,
		ImageName:    imageName,
		Tag:          tag,
	})
	if err != nil {
		return nil, err
	}

	if err := qtx.DeleteEnvVarsByDeploymentId(ctx, deploymentId); err != nil {
		return nil, err
	}

	vars := make([]*entity.EnvVar, len(envs))
	for i, e := range envs {
		variable, err := qtx.CreateDeploymentEnvVar(ctx, sqlc.CreateDeploymentEnvVarParams{
			DeploymentID: deploymentId,
			Name:         e.Name,
			Value:        e.Value,
		})
		if err != nil {
			return nil, err
		}

		vars[i] = &entity.EnvVar{
			Name:  variable.Name,
			Value: variable.Value,
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.ImageDeployment{
		Id:            d.DeploymentID,
		Status:        string(d.Status),
		ServiceName:   d.ServiceName,
		ImageName:     d.ImageName,
		Tag:           d.ImageTag,
		Port:          d.Port,
		EnvironmentId: d.EnvironmentID,
		EnvVars:       vars,
	}, nil
}

func (dr *DeploymentRepository) CreateIngressDeployment(
	ctx context.Context,
	serviceName string,
	port string,
	environmentId int64,
	hosts []*value.IngressHost,
) (*entity.IngressDeployment, error) {
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := dr.queries.WithTx(tx)
	d, err := qtx.CreateIngressDeployment(ctx, sqlc.CreateIngressDeploymentParams{
		Name:          serviceName,
		Port:          port,
		EnvironmentID: environmentId,
	})
	if err != nil {
		return nil, err
	}

	for _, h := range hosts {
		createdHost, err := qtx.CreateIngressHost(ctx, sqlc.CreateIngressHostParams{
			DeploymentID: d.DeploymentID,
			Host:         h.Host,
		})
		if err != nil {
			return nil, err
		}

		for _, p := range h.Paths {
			deployment, err := qtx.GetEnvironmentDeploymentByName(ctx, sqlc.GetEnvironmentDeploymentByNameParams{
				Name:          p.ServiceName,
				EnvironmentID: environmentId,
			})
			if err != nil {
				return nil, err
			}

			_, err = qtx.CreateIngressPath(ctx, sqlc.CreateIngressPathParams{
				IngressHostID: createdHost.ID,
				DeploymentID:  deployment.ID,
				Path:          p.Path,
				PathType:      string(p.PathType),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create ingress path: %w", err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.IngressDeployment{
		Id:            d.DeploymentID,
		Name:          d.DeploymentName,
		Port:          d.DeploymentPort,
		EnvironmentId: d.DeploymentEnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) UpdateIngressDeployment(
	ctx context.Context,
	deploymentId int64,
	port string,
	environmentId int64,
	hosts []*value.IngressHost,
) (*entity.IngressDeployment, error) {
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := dr.queries.WithTx(tx)

	d, err := qtx.UpdateIngressDeployment(ctx, sqlc.UpdateIngressDeploymentParams{
		Port:         port,
		DeploymentID: deploymentId,
	})
	if err != nil {
		return nil, err
	}

	if err := qtx.DeleteIngressPathsByDeploymentId(ctx, deploymentId); err != nil {
		return nil, err
	}
	if err := qtx.DeleteIngressHostsByDeploymentId(ctx, deploymentId); err != nil {
		return nil, err
	}

	for _, h := range hosts {
		createdHost, err := qtx.CreateIngressHost(ctx, sqlc.CreateIngressHostParams{
			DeploymentID: d.DeploymentID,
			Host:         h.Host,
		})
		if err != nil {
			return nil, err
		}

		for _, p := range h.Paths {
			deployment, err := qtx.GetEnvironmentDeploymentByName(ctx, sqlc.GetEnvironmentDeploymentByNameParams{
				Name:          p.ServiceName,
				EnvironmentID: environmentId,
			})
			if err != nil {
				return nil, err
			}

			if _, err := qtx.CreateIngressPath(ctx, sqlc.CreateIngressPathParams{
				IngressHostID: createdHost.ID,
				DeploymentID:  deployment.ID,
				Path:          p.Path,
				PathType:      string(p.PathType),
			}); err != nil {
				return nil, fmt.Errorf("failed to create ingress path: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.IngressDeployment{
		Id:            d.DeploymentID,
		Name:          d.DeploymentName,
		Port:          d.DeploymentPort,
		EnvironmentId: d.DeploymentEnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) CreateDatabaseDeployment(
	ctx context.Context,
	serviceName string,
	port string,
	environmentId int64,
) (deployment *entity.DatabaseDeployment, err error) {
	d, err := dr.queries.CreateDatabaseDeployment(ctx, sqlc.CreateDatabaseDeploymentParams{
		Name:          serviceName,
		Port:          port,
		EnvironmentID: environmentId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.DatabaseDeployment{
		Id:            d.DeploymentID,
		ServiceName:   d.Name,
		EnvironmentId: d.EnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) UpdateDatabaseDeploymentCredentials(
	ctx context.Context,
	dbName string,
	deploymentId int64,
	username string,
	password string,
) error {
	return dr.queries.UpdateDatabaseDeploymentCredentials(ctx, sqlc.UpdateDatabaseDeploymentCredentialsParams{
		Database:     utils.NullStringFromPtr(&dbName),
		Username:     utils.NullStringFromPtr(&username),
		Password:     utils.NullStringFromPtr(&password),
		DeploymentID: deploymentId,
	})
}

func (dr *DeploymentRepository) GetUserDeployment(ctx context.Context, userId int64, deploymentId int64) (*entity.Deployment, error) {
	res, err := dr.queries.GetUserDeployment(ctx, sqlc.GetUserDeploymentParams{
		DeploymentID: deploymentId,
		UserID:       userId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Deployment{
		Id:            res.ID,
		Name:          res.Name,
		Port:          res.Port,
		EnvironmentId: res.EnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) GetDeploymentWithNamespace(ctx context.Context, deploymentId int64) (*entity.Deployment, error) {
	deployment, err := dr.queries.GetDeploymentWithNamespace(ctx, deploymentId)
	if err != nil {
		return nil, err
	}
	return &entity.Deployment{
		Id:            deployment.ID,
		Name:          deployment.Name,
		Port:          deployment.Port,
		Namespace:     deployment.Namespace,
		EnvironmentId: deployment.EnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) GetDeploymentEnvs(ctx context.Context, deploymentId int64) ([]*entity.EnvVar, error) {
	envVars, err := dr.queries.GetDeploymentEnvironmentVars(ctx, deploymentId)
	if err != nil {
		return nil, err
	}

	variables := make([]*entity.EnvVar, len(envVars))
	for j, e := range envVars {
		variables[j] = &entity.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		}
	}
	return variables, nil
}

func (dr *DeploymentRepository) GetDeploymentArgs(ctx context.Context, deploymentId int64) ([]*entity.Arg, error) {
	args, err := dr.queries.GetDeploymentArgs(ctx, deploymentId)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Arg, len(args))
	for i, a := range args {
		result[i] = &entity.Arg{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return result, nil
}

func (dr *DeploymentRepository) GetDeploymentCluster(ctx context.Context, deploymentId int64) (*entity.Cluster, error) {
	res, err := dr.queries.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             res.ID,
		Name:           res.Name,
		Status:         entity.ClusterStatus(res.Status),
		IPv4Address:    utils.PtrFromNullString(res.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(res.PublicKey),
		PrivateKey:     utils.PtrFromNullString(res.PrivateKey),
		ProvisioningId: utils.PtrFromNullString(res.ProvisioningID),
		Kubeconfig:     utils.PtrFromNullString(res.Kubeconfig),
		OrganizationId: res.OrganizationID,
	}, nil
}

func (dr *DeploymentRepository) SoftDeleteDeploymentVolume(ctx context.Context, deploymentId int64) error {
	return dr.queries.SoftDeleteDeploymentVolume(ctx, sql.NullInt64{Int64: deploymentId, Valid: true})
}

func (dr *DeploymentRepository) DeleteDeployment(ctx context.Context, deploymentId int64) error {
	return dr.queries.DeleteDeployment(ctx, deploymentId)
}

func (dr *DeploymentRepository) GetAllDeploymentsWithKubeconfig(ctx context.Context) ([]*entity.DeploymentWithKubeconfig, error) {
	rows, err := dr.queries.GetDeploymentsWithKubeconfig(ctx)
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.DeploymentWithKubeconfig, len(rows))
	for i, d := range rows {
		deployments[i] = &entity.DeploymentWithKubeconfig{
			Deployment: entity.Deployment{
				Id:            d.ID,
				Name:          d.Name,
				Port:          d.Port,
				EnvironmentId: d.EnvironmentID,
				Namespace:     d.Namespace,
			},
			Kubeconfig: utils.PtrFromNullString(d.Kubeconfig),
		}
	}
	return deployments, nil
}

func (dr *DeploymentRepository) UpdateDeploymentStatus(ctx context.Context, deploymentId int64, status string) error {
	return dr.queries.UpdateDeploymentStatus(ctx, sqlc.UpdateDeploymentStatusParams{
		Status: sqlc.DeploymentStatus(status),
		ID:     deploymentId,
	})
}

func (dr *DeploymentRepository) GetEnvironmentDeploymentByName(ctx context.Context, environmentId int64, serviceName string) (*entity.Deployment, error) {
	row, err := dr.queries.GetEnvironmentDeploymentByName(ctx, sqlc.GetEnvironmentDeploymentByNameParams{
		Name:          serviceName,
		EnvironmentID: environmentId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &entity.Deployment{
		Id:            row.ID,
		Name:          row.Name,
		Port:          row.Port,
		EnvironmentId: row.EnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) GetIngressHostByName(ctx context.Context, hostName string) (*entity.IngressHostDeployment, error) {
	row, err := dr.queries.GetIngressHostByName(ctx, hostName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &entity.IngressHostDeployment{
		Host:         row.Host,
		DeploymentId: row.DeploymentID,
	}, nil
}

func (dr *DeploymentRepository) GetGitDeploymentsByRepositoryUrl(ctx context.Context, repositoryUrl string) ([]*entity.GitDeployment, error) {
	rows, err := dr.queries.GetGitDeploymentsByRepositoryUrl(ctx, repositoryUrl)
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.GitDeployment, len(rows))
	for i, d := range rows {
		envs, err := dr.queries.GetDeploymentEnvironmentVars(ctx, d.DeploymentID)
		if err != nil {
			return nil, err
		}

		envVars := make([]*entity.EnvVar, len(envs))
		for j, e := range envs {
			envVars[j] = &entity.EnvVar{
				Name:  e.Name,
				Value: e.Value,
			}
		}

		args, err := dr.queries.GetDeploymentArgs(ctx, d.DeploymentID)
		if err != nil {
			return nil, err
		}

		deploymentArgs := make([]*entity.Arg, len(args))
		for j, a := range args {
			deploymentArgs[j] = &entity.Arg{
				Name:  a.Name,
				Value: a.Value,
			}
		}

		deployments[i] = &entity.GitDeployment{
			Id:                    d.DeploymentID,
			Name:                  d.Name,
			Status:                string(d.Status),
			Port:                  d.Port,
			EnvironmentId:         d.EnvironmentID,
			GitUrl:                d.Url,
			ProjectRepositoryPath: d.ProjectPath,
			DockerfilePath:        d.DockerfilePath,
			EnvVars:               envVars,
			Args:                  deploymentArgs,
		}
	}
	return deployments, nil
}

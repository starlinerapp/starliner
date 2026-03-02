package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type EnvironmentRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.EnvironmentRepository = (*EnvironmentRepository)(nil)

func NewEnvironmentRepository(queries *sqlc.Queries) interfaces.EnvironmentRepository {
	return &EnvironmentRepository{queries: queries}
}

func (er *EnvironmentRepository) CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*entity.Environment, error) {
	env, err := er.queries.CreateEnvironment(ctx, sqlc.CreateEnvironmentParams{
		Name:      name,
		Slug:      slug,
		ProjectID: projectId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Environment{
		Id:   env.ID,
		Slug: env.Slug,
		Name: env.Name,
	}, nil
}

func (er *EnvironmentRepository) GetEnvironmentAuthorizedUsers(ctx context.Context, clusterId int64) (users []int64, err error) {
	users, err = er.queries.GetEnvironmentAuthorizedUsers(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (er *EnvironmentRepository) GetEnvironmentCluster(ctx context.Context, environmentId int64) (*entity.Cluster, error) {
	cluster, err := er.queries.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		Status:         entity.ClusterStatus(cluster.Status),
		IPv4Address:    utils.PtrFromNullString(cluster.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(cluster.PublicKey),
		PrivateKey:     utils.PtrFromNullString(cluster.PrivateKey),
		ProvisioningId: utils.PtrFromNullString(cluster.ProvisioningID),
		Kubeconfig:     utils.PtrFromNullString(cluster.Kubeconfig),
		OrganizationId: cluster.OrganizationID,
	}, nil
}

func (er *EnvironmentRepository) GetEnvironmentImageDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.ImageDeployment, error) {
	rows, err := er.queries.GetEnvironmentImageDeployments(ctx, sqlc.GetEnvironmentImageDeploymentsParams{
		EnvironmentID: environmentId,
		ID:            userId,
	})
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.ImageDeployment, len(rows))
	for i, d := range rows {
		envVars, err := er.queries.GetImageEnvironmentVars(ctx, d.DeploymentID)
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

		deployments[i] = &entity.ImageDeployment{
			Id:            d.DeploymentID,
			Status:        utils.PtrFromNullString(d.Status),
			ServiceName:   d.ServiceName,
			ImageName:     d.ImageName,
			Tag:           d.Tag,
			Port:          d.Port,
			EnvironmentId: d.EnvironmentID,
			EnvVars:       variables,
		}
	}
	return deployments, nil
}

func (er *EnvironmentRepository) GetEnvironmentIngressDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.IngressDeployment, error) {
	rows, err := er.queries.GetEnvironmentIngressDeployments(ctx, sqlc.GetEnvironmentIngressDeploymentsParams{
		EnvironmentID: environmentId,
		ID:            userId,
	})
	if err != nil {
		return nil, err
	}

	depByID := map[int64]*entity.IngressDeployment{}
	hostByDep := map[int64]map[int64]*entity.IngressHost{}

	for _, r := range rows {
		dep, exists := depByID[r.DeploymentID]
		if !exists {
			dep = &entity.IngressDeployment{
				Id:            r.DeploymentID,
				EnvironmentId: r.EnvironmentID,
				Status:        utils.PtrFromNullString(r.Status),
				Name:          r.DeploymentName,
				Port:          r.Port,
				IngressHosts:  []*entity.IngressHost{},
			}
			depByID[r.DeploymentID] = dep
			hostByDep[r.DeploymentID] = map[int64]*entity.IngressHost{}
		}

		if !r.HostID.Valid {
			continue
		}

		hID := r.HostID.Int64
		hostMap := hostByDep[r.DeploymentID]

		host, exists := hostMap[hID]
		if !exists {
			host = &entity.IngressHost{
				Host:  r.Host.String,
				Paths: []*entity.IngressPath{},
			}
			hostMap[hID] = host
			dep.IngressHosts = append(dep.IngressHosts, host)
		}

		if !r.PathID.Valid {
			continue
		}

		serviceName := ""
		if r.ServiceName.Valid {
			serviceName = r.ServiceName.String
		}

		path := ""
		if r.Path.Valid {
			path = r.Path.String
		}

		pathType := ""
		if r.PathType.Valid {
			pathType = r.PathType.String
		}

		host.Paths = append(host.Paths, &entity.IngressPath{
			Path:        path,
			PathType:    entity.PathType(pathType),
			ServiceName: serviceName,
		})
	}

	out := make([]*entity.IngressDeployment, 0, len(depByID))
	for _, dep := range depByID {
		out = append(out, dep)
	}
	return out, nil
}

func (er *EnvironmentRepository) GetEnvironmentDatabaseDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.DatabaseDeployment, error) {
	rows, err := er.queries.GetEnvironmentDatabaseDeployments(ctx, sqlc.GetEnvironmentDatabaseDeploymentsParams{
		EnvironmentID: environmentId,
		ID:            userId,
	})
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.DatabaseDeployment, len(rows))
	for i, d := range rows {
		deployments[i] = &entity.DatabaseDeployment{
			Id:            d.DeploymentID,
			Name:          d.Name,
			Status:        utils.PtrFromNullString(d.Status),
			Username:      d.Username,
			Password:      d.Password,
			Port:          d.Port,
			EnvironmentId: d.EnvironmentID,
		}
	}
	return deployments, nil
}

func (er *EnvironmentRepository) GetEnvironmentDeploymentByName(ctx context.Context, name string, environmentId int64) (*entity.Deployment, error) {
	d, err := er.queries.GetEnvironmentDeploymentByName(ctx, sqlc.GetEnvironmentDeploymentByNameParams{
		Name:          name,
		EnvironmentID: environmentId,
	})

	if err != nil {
		return nil, err
	}

	return &entity.Deployment{
		Id:            d.ID,
		Name:          d.Name,
		Port:          d.Port,
		EnvironmentId: d.EnvironmentID,
	}, nil
}

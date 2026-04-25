package repository

import (
	"context"
	"database/sql"
	"errors"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type EnvironmentRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

var _ interfaces.EnvironmentRepository = (*EnvironmentRepository)(nil)

func NewEnvironmentRepository(db *sql.DB, queries *sqlc.Queries) interfaces.EnvironmentRepository {
	return &EnvironmentRepository{
		db:      db,
		queries: queries,
	}
}

func (er *EnvironmentRepository) CreateEnvironment(ctx context.Context, name string, namespace string, slug string, projectId int64) (*entity.Environment, error) {
	env, err := er.queries.CreateEnvironment(ctx, sqlc.CreateEnvironmentParams{
		Name:      name,
		Namespace: namespace,
		Slug:      slug,
		ProjectID: projectId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Environment{
		Id:        env.ID,
		Slug:      env.Slug,
		Namespace: env.Namespace,
		Name:      env.Name,
	}, nil
}

func (er *EnvironmentRepository) DeleteEnvironment(ctx context.Context, environmentId int64) error {
	return er.queries.DeleteEnvironment(ctx, environmentId)
}

func (er *EnvironmentRepository) DuplicateEnvironment(
	ctx context.Context,
	name string,
	namespace string,
	slug string,
	projectId int64,
	sourceEnvironmentId int64,
	uniqueIdentifier string,
	connectedBranch *string,
) (*entity.Environment, error) {
	tx, err := er.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	qtx := er.queries.WithTx(tx)

	var newEnv sqlc.Environment
	if connectedBranch != nil {
		newEnv, err = qtx.CreateEnvironmentWithConnectedBranch(ctx, sqlc.CreateEnvironmentWithConnectedBranchParams{
			Name:            name,
			Slug:            slug,
			Namespace:       namespace,
			ProjectID:       projectId,
			ConnectedBranch: *connectedBranch,
		})
	} else {
		newEnv, err = qtx.CreateEnvironment(ctx, sqlc.CreateEnvironmentParams{
			Name: name, Slug: slug, Namespace: namespace, ProjectID: projectId,
		})
	}
	if err != nil {
		return nil, err
	}

	if err := er.copyDeployments(ctx, qtx, sourceEnvironmentId, newEnv.ID, uniqueIdentifier); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &entity.Environment{
		Id:        newEnv.ID,
		Slug:      newEnv.Slug,
		Name:      newEnv.Name,
		Namespace: newEnv.Namespace,
	}, nil
}

func (er *EnvironmentRepository) CreatePreviewEnvironment(
	ctx context.Context,
	name string,
	namespace string,
	slug string,
	projectId int64,
	sourceEnvironmentId int64,
	uniqueIdentifier string,
	connectedBranch *string,
	githubRepositoryId int64,
	prNumber int,
) (*entity.Environment, error) {
	tx, err := er.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	qtx := er.queries.WithTx(tx)

	newEnv, err := qtx.CreatePreviewEnvironment(ctx, sqlc.CreatePreviewEnvironmentParams{
		Name:               name,
		Slug:               slug,
		Namespace:          namespace,
		ProjectID:          projectId,
		ConnectedBranch:    *connectedBranch,
		GithubRepositoryID: githubRepositoryId,
		PrNumber:           int64(prNumber),
	})
	if err != nil {
		return nil, err
	}

	if err := er.copyDeployments(ctx, qtx, sourceEnvironmentId, newEnv.ID, uniqueIdentifier); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &entity.Environment{
		Id:        newEnv.ID,
		Slug:      newEnv.Slug,
		Name:      newEnv.Name,
		Namespace: newEnv.Namespace,
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

func (er *EnvironmentRepository) GetEnvironmentById(ctx context.Context, environmentId int64) (*entity.Environment, error) {
	env, err := er.queries.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	return &entity.Environment{
		Id:              env.ID,
		Slug:            env.Slug,
		Name:            env.Name,
		Namespace:       env.Namespace,
		ConnectedBranch: env.ConnectedBranch,
	}, nil
}

func (er *EnvironmentRepository) GetUserEnvironmentGitDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.GitDeployment, error) {
	rows, err := er.queries.GetUserEnvironmentGitDeployments(ctx, sqlc.GetUserEnvironmentGitDeploymentsParams{
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
		UserID:        userId,
	})
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.GitDeployment, len(rows))
	for i, r := range rows {
		envVars, err := er.queries.GetDeploymentEnvironmentVars(ctx, r.DeploymentID)
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

		args, err := er.queries.GetGitDeploymentArgs(ctx, r.DeploymentID)
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
			Id:                    r.DeploymentID,
			Name:                  r.Name,
			Port:                  r.Port,
			Status:                string(r.Status),
			EnvironmentId:         utils.PtrFromNullInt64(r.EnvironmentID),
			GitUrl:                r.Url,
			ProjectRepositoryPath: r.ProjectPath,
			DockerfilePath:        r.DockerfilePath,
			EnvVars:               variables,
			Args:                  deploymentArgs,
		}
	}

	return deployments, nil
}

func (er *EnvironmentRepository) GetEnvironmentGitDeployments(ctx context.Context, environmentId int64) ([]*entity.GitDeployment, error) {
	rows, err := er.queries.GetEnvironmentGitDeployments(ctx, utils.NullInt64FromPtr(&environmentId))
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.GitDeployment, len(rows))
	for i, r := range rows {
		envVars, err := er.queries.GetDeploymentEnvironmentVars(ctx, r.DeploymentID)
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

		args, err := er.queries.GetGitDeploymentArgs(ctx, r.DeploymentID)
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
			Id:                    r.DeploymentID,
			Name:                  r.Name,
			Port:                  r.Port,
			Status:                string(r.Status),
			EnvironmentId:         utils.PtrFromNullInt64(r.EnvironmentID),
			GitUrl:                r.Url,
			ProjectRepositoryPath: r.ProjectPath,
			DockerfilePath:        r.DockerfilePath,
			EnvVars:               variables,
			Args:                  deploymentArgs,
		}
	}

	return deployments, nil
}

func (er *EnvironmentRepository) GetUserEnvironmentImageDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.ImageDeployment, error) {
	rows, err := er.queries.GetUserEnvironmentImageDeployments(ctx, sqlc.GetUserEnvironmentImageDeploymentsParams{
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
		UserID:        userId,
	})
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.ImageDeployment, len(rows))
	for i, d := range rows {
		envVars, err := er.queries.GetDeploymentEnvironmentVars(ctx, d.DeploymentID)
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
			Id:              d.DeploymentID,
			Status:          string(d.Status),
			ServiceName:     d.ServiceName,
			ImageName:       d.ImageName,
			Tag:             d.Tag,
			Port:            d.Port,
			EnvironmentId:   utils.PtrFromNullInt64(d.EnvironmentID),
			VolumeSizeMiB:   utils.PtrFromNullInt32(d.VolumeSizeMib),
			VolumeMountPath: utils.PtrFromNullString(d.MountPath),
			EnvVars:         variables,
		}
	}
	return deployments, nil
}

func (er *EnvironmentRepository) GetEnvironmentImageDeployments(ctx context.Context, environmentId int64) ([]*entity.ImageDeployment, error) {
	rows, err := er.queries.GetEnvironmentImageDeployments(ctx, utils.NullInt64FromPtr(&environmentId))
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.ImageDeployment, len(rows))
	for i, d := range rows {
		envVars, err := er.queries.GetDeploymentEnvironmentVars(ctx, d.DeploymentID)
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
			Id:              d.DeploymentID,
			Status:          string(d.Status),
			ServiceName:     d.ServiceName,
			ImageName:       d.ImageName,
			Tag:             d.Tag,
			Port:            d.Port,
			EnvironmentId:   utils.PtrFromNullInt64(d.EnvironmentID),
			VolumeSizeMiB:   utils.PtrFromNullInt32(d.VolumeSizeMib),
			VolumeMountPath: utils.PtrFromNullString(d.MountPath),
			EnvVars:         variables,
		}
	}
	return deployments, nil
}

func (er *EnvironmentRepository) GetUserEnvironmentIngressDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.IngressDeployment, error) {
	rows, err := er.queries.GetUserEnvironmentIngressDeployments(ctx, sqlc.GetUserEnvironmentIngressDeploymentsParams{
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
		UserID:        userId,
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
				EnvironmentId: utils.PtrFromNullInt64(r.EnvironmentID),
				Status:        string(r.Status),
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

func (er *EnvironmentRepository) GetEnvironmentIngressDeployments(ctx context.Context, environmentId int64) ([]*entity.IngressDeployment, error) {
	rows, err := er.queries.GetEnvironmentIngressDeployments(ctx, utils.NullInt64FromPtr(&environmentId))
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
				EnvironmentId: utils.PtrFromNullInt64(r.EnvironmentID),
				Status:        string(r.Status),
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

func (er *EnvironmentRepository) GetEnvironmentIngressDeploymentByName(ctx context.Context, environmentId int64, name string) (*entity.IngressDeployment, error) {
	rows, err := er.queries.GetEnvironmentIngressDeploymentsByName(ctx, sqlc.GetEnvironmentIngressDeploymentsByNameParams{
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
		Name:          name,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	first := rows[0]

	dep := &entity.IngressDeployment{
		Id:            first.DeploymentID,
		EnvironmentId: utils.PtrFromNullInt64(first.EnvironmentID),
		Status:        string(first.Status),
		Name:          first.DeploymentName,
		Port:          first.Port,
		IngressHosts:  []*entity.IngressHost{},
	}

	var host *entity.IngressHost
	if first.HostID.Valid {
		host = &entity.IngressHost{
			Host:  first.Host.String,
			Paths: []*entity.IngressPath{},
		}
		dep.IngressHosts = append(dep.IngressHosts, host)
	}

	for _, r := range rows {
		if !r.PathID.Valid || host == nil {
			continue
		}

		host.Paths = append(host.Paths, &entity.IngressPath{
			Path:        r.Path.String,
			PathType:    entity.PathType(r.PathType.String),
			ServiceName: r.ServiceName.String,
		})
	}

	return dep, nil
}

func (er *EnvironmentRepository) GetUserEnvironmentDatabaseDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.DatabaseDeployment, error) {
	rows, err := er.queries.GetUserEnvironmentDatabaseDeployments(ctx, sqlc.GetUserEnvironmentDatabaseDeploymentsParams{
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
		UserID:        userId,
	})
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.DatabaseDeployment, len(rows))
	for i, d := range rows {
		deployments[i] = &entity.DatabaseDeployment{
			Id:            d.DeploymentID,
			ServiceName:   d.Name,
			Status:        string(d.Status),
			Database:      utils.PtrFromNullString(d.Database),
			Username:      utils.PtrFromNullString(d.Username),
			Password:      utils.PtrFromNullString(d.Password),
			Port:          d.Port,
			EnvironmentId: utils.PtrFromNullInt64(d.EnvironmentID),
		}
	}
	return deployments, nil
}

func (er *EnvironmentRepository) GetEnvironmentDatabaseDeployments(ctx context.Context, environmentId int64) ([]*entity.DatabaseDeployment, error) {
	rows, err := er.queries.GetEnvironmentDatabaseDeployments(ctx, utils.NullInt64FromPtr(&environmentId))
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.DatabaseDeployment, len(rows))
	for i, d := range rows {
		deployments[i] = &entity.DatabaseDeployment{
			Id:            d.DeploymentID,
			ServiceName:   d.Name,
			Status:        string(d.Status),
			Database:      utils.PtrFromNullString(d.Database),
			Username:      utils.PtrFromNullString(d.Username),
			Password:      utils.PtrFromNullString(d.Password),
			Port:          d.Port,
			EnvironmentId: utils.PtrFromNullInt64(d.EnvironmentID),
		}
	}
	return deployments, nil
}

func (er *EnvironmentRepository) GetEnvironmentDeploymentByName(ctx context.Context, name string, environmentId int64) (*entity.Deployment, error) {
	d, err := er.queries.GetEnvironmentDeploymentByName(ctx, sqlc.GetEnvironmentDeploymentByNameParams{
		Name:          name,
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
	})

	if err != nil {
		return nil, err
	}

	return &entity.Deployment{
		Id:            d.ID,
		Name:          d.Name,
		Port:          d.Port,
		EnvironmentId: utils.PtrFromNullInt64(d.EnvironmentID),
	}, nil
}

func (er *EnvironmentRepository) GetEnvironmentGitDeploymentBuilds(ctx context.Context, environmentId int64) ([]*entity.GitDeploymentBuild, error) {
	rows, err := er.queries.GetEnvironmentGitDeploymentBuilds(ctx, utils.NullInt64FromPtr(&environmentId))
	if err != nil {
		return nil, err
	}

	builds := make([]*entity.GitDeploymentBuild, len(rows))
	for i, row := range rows {
		args, err := er.queries.GetGitDeploymentArgs(ctx, row.DeploymentID)
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

		builds[i] = &entity.GitDeploymentBuild{
			BuildId:        row.BuildID,
			DeploymentId:   row.DeploymentID,
			DeploymentName: row.DeploymentName,
			CommitHash:     utils.PtrFromNullString(row.CommitHash),
			Source:         row.Source,
			Status:         entity.BuildStatus(row.Status),
			GitUrl:         row.Url,
			ProjectPath:    row.ProjectPath,
			DockerfilePath: row.DockerfilePath,
			CreatedAt:      row.CreatedAt,
			Args:           deploymentArgs,
		}
	}

	return builds, nil
}

func (er *EnvironmentRepository) GetEnvironmentBranch(ctx context.Context, environmentId int64) (string, error) {
	branch, err := er.queries.GetEnvironmentBranch(ctx, environmentId)
	if err != nil {
		return "", err
	}

	return branch, nil
}

func (er *EnvironmentRepository) UpdateEnvironmentBranch(ctx context.Context, environmentId int64, branch string) error {
	return er.queries.UpdateEnvironmentBranch(ctx, sqlc.UpdateEnvironmentBranchParams{
		ConnectedBranch: branch,
		ID:              environmentId,
	})
}

func (er *EnvironmentRepository) GetPreviewEnvironment(ctx context.Context, gitHubRepositoryId int64, prNumber int) (*entity.PreviewEnvironment, error) {
	env, err := er.queries.GetPreviewEnvironment(ctx, sqlc.GetPreviewEnvironmentParams{
		GithubRepositoryID: gitHubRepositoryId,
		PrNumber:           int64(prNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.PreviewEnvironment{
		Id:                 env.ID,
		Slug:               env.Slug,
		Name:               env.Name,
		Namespace:          env.Namespace,
		GithubRepositoryId: env.GithubRepositoryID,
		PrNumber:           env.PrNumber,
	}, nil
}

func (er *EnvironmentRepository) GetEnvironmentProject(ctx context.Context, environmentId int64) (*entity.Project, error) {
	project, err := er.queries.GetEnvironmentProject(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	return &entity.Project{
		Id:                    project.ID,
		Name:                  project.Name,
		TeamId:                project.TeamID,
		PrEnvironmentsEnabled: utils.BoolPtrFromNullBool(project.PreviewEnvironmentsEnabled),
		ClusterId:             utils.PtrFromNullInt64(project.ClusterID),
		CreatedAt:             project.CreatedAt,
	}, nil
}

func (er *EnvironmentRepository) copyDeployments(
	ctx context.Context,
	qtx *sqlc.Queries,
	sourceEnvironmentId int64,
	newEnvID int64,
	uniqueIdentifier string,
) error {
	databaseDeployments, err := qtx.GetEnvironmentDatabaseDeployments(ctx, utils.NullInt64FromPtr(&sourceEnvironmentId))
	if err != nil {
		return err
	}
	for _, d := range databaseDeployments {
		if _, err := qtx.CreateDatabaseDeployment(ctx, sqlc.CreateDatabaseDeploymentParams{
			Name:          d.Name,
			Port:          d.Port,
			EnvironmentID: utils.NullInt64FromPtr(&newEnvID),
		}); err != nil {
			return err
		}
	}

	imageDeployments, err := qtx.GetEnvironmentImageDeployments(ctx, utils.NullInt64FromPtr(&sourceEnvironmentId))
	if err != nil {
		return err
	}
	for _, d := range imageDeployments {
		imageDeployment, err := qtx.CreateImageDeployment(ctx, sqlc.CreateImageDeploymentParams{
			ServiceName:   d.ServiceName,
			Port:          d.Port,
			EnvironmentID: utils.NullInt64FromPtr(&newEnvID),
			ImageName:     d.ImageName,
			Tag:           d.Tag,
		})
		if err != nil {
			return err
		}
		if err := er.copyEnvVars(ctx, qtx, d.DeploymentID, imageDeployment.DeploymentID); err != nil {
			return err
		}
		volume, err := qtx.GetDeploymentVolume(ctx, utils.NullInt64FromPtr(&d.DeploymentID))
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			if _, err = qtx.CreateDeploymentVolume(ctx, sqlc.CreateDeploymentVolumeParams{
				DeploymentID:  utils.NullInt64FromPtr(&imageDeployment.DeploymentID),
				VolumeSizeMib: volume.VolumeSizeMib,
				MountPath:     volume.MountPath,
			}); err != nil {
				return err
			}
		}
	}

	gitDeployments, err := qtx.GetEnvironmentGitDeployments(ctx, utils.NullInt64FromPtr(&sourceEnvironmentId))
	if err != nil {
		return err
	}
	for _, d := range gitDeployments {
		newDeployment, err := qtx.CreateGitDeployment(ctx, sqlc.CreateGitDeploymentParams{
			Name:           d.Name,
			Port:           d.Port,
			EnvironmentID:  utils.NullInt64FromPtr(&newEnvID),
			Url:            d.Url,
			ProjectPath:    d.ProjectPath,
			DockerfilePath: d.DockerfilePath,
		})
		if err != nil {
			return err
		}
		if err := er.copyEnvVars(ctx, qtx, d.DeploymentID, newDeployment.DeploymentID); err != nil {
			return err
		}
	}

	return er.copyIngressDeployments(ctx, qtx, sourceEnvironmentId, newEnvID, uniqueIdentifier)
}

func (er *EnvironmentRepository) copyEnvVars(
	ctx context.Context,
	qtx *sqlc.Queries,
	sourceDeploymentID int64,
	targetDeploymentID int64,
) error {
	envVars, err := qtx.GetDeploymentEnvironmentVars(ctx, sourceDeploymentID)
	if err != nil {
		return err
	}
	for _, e := range envVars {
		if _, err := qtx.CreateDeploymentEnvVar(ctx, sqlc.CreateDeploymentEnvVarParams{
			DeploymentID: targetDeploymentID,
			Name:         e.Name,
			Value:        e.Value,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (er *EnvironmentRepository) copyIngressDeployments(
	ctx context.Context,
	qtx *sqlc.Queries,
	sourceEnvironmentId int64,
	newEnvID int64,
	uniqueIdentifier string,
) error {
	ingressRows, err := qtx.GetEnvironmentIngressDeployments(ctx, utils.NullInt64FromPtr(&sourceEnvironmentId))
	if err != nil {
		return err
	}

	depByID := map[int64]*entity.IngressDeployment{}
	hostByDep := map[int64]map[int64]*entity.IngressHost{}

	for _, r := range ingressRows {
		dep, exists := depByID[r.DeploymentID]
		if !exists {
			dep = &entity.IngressDeployment{
				Id:            r.DeploymentID,
				EnvironmentId: utils.PtrFromNullInt64(r.EnvironmentID),
				Status:        string(r.Status),
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
			hostValue := ""
			if r.Host.Valid {
				hostValue = r.Host.String
			}
			host = &entity.IngressHost{
				Host:  uniqueIdentifier + "-" + hostValue,
				Paths: []*entity.IngressPath{},
			}
			hostMap[hID] = host
			dep.IngressHosts = append(dep.IngressHosts, host)
		}
		if !r.PathID.Valid {
			continue
		}
		serviceName, path, pathType := "", "", ""
		if r.ServiceName.Valid {
			serviceName = r.ServiceName.String
		}
		if r.Path.Valid {
			path = r.Path.String
		}
		if r.PathType.Valid {
			pathType = r.PathType.String
		}
		host.Paths = append(host.Paths, &entity.IngressPath{
			Path:        path,
			PathType:    entity.PathType(pathType),
			ServiceName: serviceName,
		})
	}

	for _, i := range depByID {
		createdIngress, err := qtx.CreateIngressDeployment(ctx, sqlc.CreateIngressDeploymentParams{
			Name:          i.Name,
			Port:          i.Port,
			EnvironmentID: utils.NullInt64FromPtr(&newEnvID),
		})
		if err != nil {
			return err
		}
		for _, h := range i.IngressHosts {
			createdHost, err := qtx.CreateIngressHost(ctx, sqlc.CreateIngressHostParams{
				DeploymentID: createdIngress.DeploymentID,
				Host:         h.Host,
			})
			if err != nil {
				return err
			}
			for _, p := range h.Paths {
				targetDeployment, err := qtx.GetEnvironmentDeploymentByName(ctx, sqlc.GetEnvironmentDeploymentByNameParams{
					Name:          p.ServiceName,
					EnvironmentID: utils.NullInt64FromPtr(&newEnvID),
				})
				if err != nil {
					return err
				}
				if _, err = qtx.CreateIngressPath(ctx, sqlc.CreateIngressPathParams{
					IngressHostID: createdHost.ID,
					DeploymentID:  targetDeployment.ID,
					Path:          p.Path,
					PathType:      string(p.PathType),
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

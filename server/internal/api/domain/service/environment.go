package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

type GitProvisionStrategy int

const (
	GitProvisionReuseSourceImage GitProvisionStrategy = iota
	GitProvisionBuildFromBranch
)

type CloneSpec struct {
	Name              string
	Namespace         string
	Slug              string
	ProjectID         int64
	SourceEnvironment int64
	UniquePrefix      string
	ConnectedBranch   *string
	Preview           *entity.PreviewCloneMetadata
}

type ProvisionOptions struct {
	GitStrategy            GitProvisionStrategy
	SourceEnvIDForGitReuse int64
	GitBranch              string
	PreviewURLs            *[]string
}

type EnvironmentService struct {
	cfg                   *conf.Config
	crypto                corePort.Crypto
	queue                 port.Queue
	buildService          *BuildService
	buildRepository       interfaces.BuildRepository
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	gitDeploymentService  *GitDeploymentService
	normalizer            *coreService.NormalizerService
	parser                *ParserService
	resolver              *ResolverService
}

func NewEnvironmentService(
	cfg *conf.Config,
	crypto corePort.Crypto,
	queue port.Queue,
	buildService *BuildService,
	buildRepository interfaces.BuildRepository,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	gitDeploymentService *GitDeploymentService,
	normalizer *coreService.NormalizerService,
	parser *ParserService,
	resolver *ResolverService,
) *EnvironmentService {
	return &EnvironmentService{
		cfg:                   cfg,
		crypto:                crypto,
		queue:                 queue,
		buildService:          buildService,
		buildRepository:       buildRepository,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		gitDeploymentService:  gitDeploymentService,
		normalizer:            normalizer,
		parser:                parser,
		resolver:              resolver,
	}
}

func (es *EnvironmentService) ValidateUserPermission(ctx context.Context, userId int64, environmentId int64) error {
	users, err := es.environmentRepository.GetEnvironmentAuthorizedUsers(ctx, environmentId)
	if err != nil {
		return err
	}

	found := false
	for _, user := range users {
		if user == userId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not authorized")
	}
	return nil
}

func (es *EnvironmentService) TearDownEnvironmentDeployments(ctx context.Context, env *entity.Environment) error {
	ingresses, err := es.environmentRepository.GetEnvironmentIngressDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	gitDeployments, err := es.environmentRepository.GetEnvironmentGitDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	images, err := es.environmentRepository.GetEnvironmentImageDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	databases, err := es.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	type deploymentIDAndName struct {
		id          int64
		serviceName string
	}
	var toRemove []deploymentIDAndName

	for _, d := range ingresses {
		toRemove = append(toRemove, deploymentIDAndName{d.Id, d.Name})
	}
	for _, d := range gitDeployments {
		toRemove = append(toRemove, deploymentIDAndName{d.Id, d.Name})
	}
	for _, d := range images {
		toRemove = append(toRemove, deploymentIDAndName{d.Id, d.ServiceName})
	}
	for _, d := range databases {
		toRemove = append(toRemove, deploymentIDAndName{d.Id, d.ServiceName})
	}

	for _, d := range toRemove {
		cluster, err := es.deploymentRepository.GetDeploymentCluster(ctx, d.id)
		if err != nil {
			return err
		}

		if cluster.Kubeconfig == nil {
			return fmt.Errorf("cluster kubeconfig is nil")
		}
		kubeconfigBase64, err := es.crypto.Decrypt(*cluster.Kubeconfig)
		if err != nil {
			return err
		}

		normalizedDeploymentName, err := es.normalizer.FormatToDNS1123(d.serviceName)
		if err != nil {
			return err
		}

		if err = es.deploymentRepository.SoftDeleteDeploymentVolume(ctx, d.id); err != nil {
			return err
		}

		if err = es.queue.PublishDeleteDeployment(&coreValue.Deployment{
			DeploymentId:     d.id,
			DeploymentName:   normalizedDeploymentName,
			Namespace:        env.Namespace,
			KubeconfigBase64: kubeconfigBase64,
		}); err != nil {
			log.Printf("error publishing: %v", err)
		}
	}
	return nil
}

func (es *EnvironmentService) RandomPrefix(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func (es *EnvironmentService) CloneAndProvision(
	ctx context.Context,
	spec CloneSpec,
	opts ProvisionOptions,
) (*entity.Environment, error) {
	env, err := es.environmentRepository.CloneEnvironment(
		ctx,
		spec.Name,
		spec.Namespace,
		spec.Slug,
		spec.ProjectID,
		spec.SourceEnvironment,
		spec.UniquePrefix,
		spec.ConnectedBranch,
		spec.Preview,
	)
	if err != nil {
		return nil, err
	}

	deployments, err := es.ListEnvironmentDeployments(ctx, env.Id)
	if err != nil {
		return nil, err
	}

	cluster, err := es.environmentRepository.GetEnvironmentCluster(ctx, env.Id)
	if err != nil {
		return nil, err
	}
	if cluster.IPv4Address == nil || *cluster.IPv4Address == "" {
		return nil, fmt.Errorf("cluster ipv4 address is not set")
	}
	if cluster.Kubeconfig == nil {
		return nil, fmt.Errorf("cluster kubeconfig is nil")
	}

	kubeconfigBase64, err := es.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return nil, err
	}

	if err := es.ProvisionDeployments(ctx, env, deployments, opts, kubeconfigBase64, *cluster.IPv4Address); err != nil {
		return nil, err
	}

	return env, nil
}

func (es *EnvironmentService) ListEnvironmentDeployments(
	ctx context.Context,
	environmentId int64,
) (*value.Deployments, error) {
	ingresses, err := es.environmentRepository.GetEnvironmentIngressDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	git, err := es.environmentRepository.GetEnvironmentGitDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	gitDeployments := make([]*value.GitDeployment, len(git))
	for i, d := range git {
		normalizedServiceName, err := es.normalizer.FormatToDNS1123(d.Name)
		if err != nil {
			return nil, err
		}
		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		gitDeployments[i] = value.NewGitDeployment(d, internalEndpoint)
	}

	images, err := es.environmentRepository.GetEnvironmentImageDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	imageDeployments := make([]*value.ImageDeployment, len(images))
	for i, d := range images {
		normalizedServiceName, err := es.normalizer.FormatToDNS1123(d.ServiceName)
		if err != nil {
			return nil, err
		}
		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		imageDeployments[i] = value.NewImageDeployment(d, internalEndpoint)
	}

	databases, err := es.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	databaseDeployments := make([]*value.DatabaseDeployment, len(databases))
	for i, d := range databases {
		var password *string
		if d.Password != nil {
			decrypted, err := es.crypto.Decrypt(*d.Password)
			if err != nil {
				return nil, err
			}
			password = &decrypted
		}

		normalizedServiceName, err := es.normalizer.FormatToDNS1123(d.ServiceName)
		if err != nil {
			return nil, err
		}
		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		databaseDeployments[i] = &value.DatabaseDeployment{
			Id:               d.Id,
			ServiceName:      d.ServiceName,
			InternalEndpoint: internalEndpoint,
			Status:           d.Status,
			Database:         d.Database,
			Username:         d.Username,
			Password:         password,
			Port:             d.Port,
		}
	}

	return &value.Deployments{
		Databases:      databaseDeployments,
		Images:         imageDeployments,
		Ingresses:      value.NewIngressDeployments(ingresses),
		GitDeployments: gitDeployments,
	}, nil
}

func (es *EnvironmentService) ProvisionDeployments(
	ctx context.Context,
	env *entity.Environment,
	deployments *value.Deployments,
	opts ProvisionOptions,
	kubeconfigBase64 string,
	clusterIPv4 string,
) error {
	var errs []error

	for _, d := range deployments.Ingresses {
		coreHosts, err := es.buildIngressCoreHosts(ctx, d, env.Id, opts.PreviewURLs)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err := es.buildService.CreateDeployOnlyBuild(ctx, d.Id, value.BuildSourceDuplicate); err != nil {
			log.Printf("failed to create ingress deploy build: %v", err)
			continue
		}

		if err := es.queue.PublishDeployIngress(&coreValue.IngressDeployment{
			IngressHosts:     coreHosts,
			DeploymentId:     d.Id,
			DeploymentName:   d.ServiceName,
			Namespace:        env.Namespace,
			KubeconfigBase64: kubeconfigBase64,
			ExpectedIP:       clusterIPv4,
		}); err != nil {
			log.Printf("error publishing ingress deploy: %v", err)
		}
	}

	for _, d := range deployments.Databases {
		normalizedServiceName, err := es.normalizer.FormatToDNS1123(d.ServiceName)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err := es.buildService.CreateDeployOnlyBuild(ctx, d.Id, value.BuildSourceDuplicate); err != nil {
			log.Printf("failed to create database deploy build: %v", err)
			continue
		}

		if err := es.queue.PublishDeployDatabase(&coreValue.Deployment{
			Namespace:        env.Namespace,
			DeploymentId:     d.Id,
			DeploymentName:   normalizedServiceName,
			KubeconfigBase64: kubeconfigBase64,
		}); err != nil {
			log.Printf("error publishing database deploy: %v", err)
		}
	}

	for _, d := range deployments.Images {
		deploymentPort, err := strconv.Atoi(d.Port)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		normalizedDeploymentName, err := es.normalizer.FormatToDNS1123(d.ServiceName)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err := es.buildService.CreateDeployOnlyBuild(ctx, d.Id, value.BuildSourceDuplicate); err != nil {
			log.Printf("failed to create image deploy build: %v", err)
			continue
		}

		if err := es.queue.PublishDeployImage(&coreValue.ImageDeployment{
			DeploymentId:          d.Id,
			DeploymentName:        normalizedDeploymentName,
			Namespace:             env.Namespace,
			KubeconfigBase64:      kubeconfigBase64,
			ImageRegistryUrl:      es.cfg.ImageRegistryUrl,
			ImageRegistryUsername: es.cfg.ImageRegistryUsername,
			ImageRegistryPassword: es.cfg.ImageRegistryPassword,
			ImageName:             d.ImageName,
			ImageTag:              d.Tag,
			Port:                  deploymentPort,
			VolumeSizeMiB:         d.VolumeSizeMiB,
			VolumeMountPath:       d.VolumeMountPath,
			EnvVars:               value.ToCoreEnvVars(d.EnvVars),
		}); err != nil {
			log.Printf("error publishing image deploy: %v", err)
		}
	}

	for _, d := range deployments.GitDeployments {
		switch opts.GitStrategy {
		case GitProvisionReuseSourceImage:
			if err := es.provisionClonedGitFromSourceImage(ctx, d, env, opts.SourceEnvIDForGitReuse, kubeconfigBase64); err != nil {
				errs = append(errs, err)
			}
		case GitProvisionBuildFromBranch:
			gitEntity, err := es.environmentRepository.GetEnvironmentGitDeployments(ctx, env.Id)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			var match *entity.GitDeployment
			for _, g := range gitEntity {
				if g.Name == d.ServiceName {
					match = g
					break
				}
			}
			if match == nil {
				errs = append(errs, fmt.Errorf("git deployment %s not found after clone", d.ServiceName))
				continue
			}
			args := entityArgsToValue(match.Args)
			if err := es.gitDeploymentService.TriggerBuild(ctx, match, opts.GitBranch, "push", args); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (es *EnvironmentService) buildIngressCoreHosts(
	ctx context.Context,
	d *value.IngressDeployment,
	environmentId int64,
	previewURLs *[]string,
) ([]coreValue.IngressHost, error) {
	coreHosts := make([]coreValue.IngressHost, 0, len(d.IngressHosts))
	for _, h := range d.IngressHosts {
		if previewURLs != nil && h.Host != "" {
			*previewURLs = append(*previewURLs, "https://"+h.Host)
		}

		ch := coreValue.IngressHost{Host: h.Host}
		ch.Paths = make([]coreValue.IngressPath, 0, len(h.Paths))

		for _, p := range h.Paths {
			target, err := es.environmentRepository.GetEnvironmentDeploymentByName(ctx, p.ServiceName, environmentId)
			if err != nil {
				return nil, err
			}

			targetPort, err := strconv.Atoi(target.Port)
			if err != nil {
				return nil, err
			}

			normalizedServiceName, err := es.normalizer.FormatToDNS1123(p.ServiceName)
			if err != nil {
				return nil, err
			}

			ch.Paths = append(ch.Paths, coreValue.IngressPath{
				Path:        p.Path,
				PathType:    coreValue.PathType(p.PathType),
				ServiceName: normalizedServiceName,
				ServicePort: targetPort,
			})
		}
		coreHosts = append(coreHosts, ch)
	}
	return coreHosts, nil
}

func (es *EnvironmentService) provisionClonedGitFromSourceImage(
	ctx context.Context,
	deployment *value.GitDeployment,
	env *entity.Environment,
	sourceEnvID int64,
	kubeconfigBase64 string,
) error {
	latestBuild, err := es.buildRepository.GetLatestGitDeploymentBuild(ctx, sourceEnvID, deployment.ServiceName)
	if err != nil {
		return err
	}
	if latestBuild.ImageName == nil {
		return fmt.Errorf("latest build for git deployment %s is nil", deployment.ServiceName)
	}

	coreEnvs := make([]*coreValue.EnvVar, 0, len(deployment.EnvVars))
	for _, e := range deployment.EnvVars {
		res, err := es.parser.Parse(e.Value)
		if err != nil {
			log.Printf("failed to parse env var: %v\n", err)
			continue
		}

		resolvedValue, err := es.resolver.Resolve(ctx, env.Id, res)
		if err != nil {
			log.Printf("failed to resolve env var: %v\n", err)
			continue
		}

		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: resolvedValue,
		})
	}

	normalizedDeploymentName, err := es.normalizer.FormatToDNS1123(deployment.ServiceName)
	if err != nil {
		return err
	}
	deploymentPort, err := strconv.Atoi(deployment.Port)
	if err != nil {
		return err
	}

	b, err := es.buildRepository.CreateBuild(ctx, deployment.Id, value.BuildSourceDuplicate)
	if err != nil {
		return err
	}

	if err := es.buildRepository.UpdateBuild(
		ctx,
		b.Id,
		value.BuildStatusSuccess,
		latestBuild.CommitHash,
		latestBuild.ImageName,
		"",
	); err != nil {
		return err
	}

	return es.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:          deployment.Id,
		DeploymentName:        normalizedDeploymentName,
		Namespace:             env.Namespace,
		KubeconfigBase64:      kubeconfigBase64,
		ImageRegistryUrl:      es.cfg.ImageRegistryUrl,
		ImageRegistryUsername: es.cfg.ImageRegistryUsername,
		ImageRegistryPassword: es.cfg.ImageRegistryPassword,
		ImageName:             *latestBuild.ImageName,
		ImageTag:              *latestBuild.CommitHash,
		Port:                  deploymentPort,
		EnvVars:               coreEnvs,
	})
}

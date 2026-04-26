package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

type EnvironmentService struct {
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	queue                 port.Queue
	crypto                corePort.Crypto
	normalizer            *coreService.NormalizerService
}

func NewEnvironmentService(
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	queue port.Queue,
	crypto corePort.Crypto,
	normalizer *coreService.NormalizerService,
) *EnvironmentService {
	return &EnvironmentService{
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		queue:                 queue,
		crypto:                crypto,
		normalizer:            normalizer,
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

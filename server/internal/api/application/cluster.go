package application

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/sse"
	corePort "starliner.app/internal/core/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
)

type ClusterApplication struct {
	clusterRepository      interfaces.ClusterRepository
	organizationRepository interfaces.OrganizationRepository
	teamRepository         interfaces.TeamRepository
	projectRepository      interfaces.ProjectRepository
	environmentRepository  interfaces.EnvironmentRepository
	deploymentRepository   interfaces.DeploymentRepository
	environmentService     *service.EnvironmentService
	organizationService    *service.OrganizationService
	crypto                 corePort.Crypto
	queue                  port.Queue
	grpcProvisionerClient  port.ProvisionerClient
	userNotificationHub    *sse.UserNotificationHub
}

func NewClusterApplication(
	clusterRepository interfaces.ClusterRepository,
	organizationRepository interfaces.OrganizationRepository,
	teamRepository interfaces.TeamRepository,
	projectRepository interfaces.ProjectRepository,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	environmentService *service.EnvironmentService,
	organizationService *service.OrganizationService,
	crypto corePort.Crypto,
	queue port.Queue,
	grpcProvisionerClient port.ProvisionerClient,
	userNotificationHub *sse.UserNotificationHub,
) *ClusterApplication {
	return &ClusterApplication{
		clusterRepository:      clusterRepository,
		organizationRepository: organizationRepository,
		teamRepository:         teamRepository,
		projectRepository:      projectRepository,
		environmentRepository:  environmentRepository,
		deploymentRepository:   deploymentRepository,
		environmentService:     environmentService,
		organizationService:    organizationService,
		crypto:                 crypto,
		queue:                  queue,
		grpcProvisionerClient:  grpcProvisionerClient,
		userNotificationHub:    userNotificationHub,
	}
}

func (ca *ClusterApplication) CreateCluster(ctx context.Context, userId int64, name string, serverType string, organizationId int64, teamId int64) (*value.Cluster, error) {
	if err := ca.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId); err != nil {
		return nil, err
	}

	team, err := ca.teamRepository.GetTeamById(ctx, teamId)
	if err != nil {
		return nil, err
	}

	// the org owner should be part of every team in the organization
	if team.OrganizationId != organizationId {
		return nil, errors.New("team does not belong to the specified organization")
	}

	cluster, err := ca.clusterRepository.CreateCluster(ctx, name, serverType, organizationId)
	if err != nil {
		return nil, fmt.Errorf("failed to persist cluster in database: %v", err)
	}

	if err := ca.teamRepository.AssignClusterToTeam(ctx, teamId, cluster.Id); err != nil {
		_ = ca.clusterRepository.DeleteCluster(ctx, cluster.Id)
		return nil, err
	}

	credential, err := ca.organizationRepository.GetOrganizationProvisioningCredential(ctx, organizationId, value.HetznerCredential)
	if err != nil {
		return nil, err
	}

	decrypted, err := ca.crypto.Decrypt(credential.Secret)
	if err != nil {
		return nil, err
	}

	err = ca.queue.PublishCreateCluster(&coreValue.ProvisionCluster{
		Id:                     cluster.Id,
		Name:                   cluster.Name,
		ServerType:             coreValue.ServerType(cluster.ServerType),
		OrganizationName:       strconv.FormatInt(cluster.OrganizationId, 10),
		ProvisioningCredential: decrypted,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return value.NewCluster(cluster), nil
}

func (ca *ClusterApplication) GetCluster(ctx context.Context, id int64) (*value.Cluster, error) {
	cluster, err := ca.clusterRepository.GetCluster(ctx, id)
	if err != nil {
		return nil, err
	}

	return value.NewCluster(cluster), nil
}

func (ca *ClusterApplication) GetClusterOrgOwnerId(ctx context.Context, clusterId int64) (int64, error) {
	return ca.clusterRepository.GetClusterOrgOwnerId(ctx, clusterId)
}

func (ca *ClusterApplication) GetUserCluster(ctx context.Context, userId int64, id int64) (*value.Cluster, error) {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return value.NewCluster(cluster), nil
}

func (ca *ClusterApplication) GetClusterPrivateKey(ctx context.Context, id int64, userId int64) ([]byte, error) {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if cluster.PrivateKey == nil {
		return nil, fmt.Errorf("cluster private key is not set")
	}

	// The private key was first base64 encoded and then encrypted
	decryptedPrivateKey, err := ca.crypto.Decrypt(*cluster.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(decryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	pemBytes, err := ca.crypto.EncodePrivateKeyToPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key to PEM: %v", err)
	}
	return pemBytes, nil
}

func (ca *ClusterApplication) DeleteCluster(ctx context.Context, userId int64, clusterId int64) error {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, clusterId)
	if err != nil {
		return err
	}

	err = ca.clusterRepository.UpdateClusterStatus(ctx, clusterId, entity.ClusterDeleted)
	if err != nil {
		return err
	}

	credential, err := ca.organizationRepository.GetOrganizationProvisioningCredential(ctx, cluster.OrganizationId, value.HetznerCredential)
	if err != nil {
		return err
	}

	decrypted, err := ca.crypto.Decrypt(credential.Secret)
	if err != nil {
		return err
	}

	err = ca.queue.PublishDeleteCluster(&coreValue.DeleteCluster{
		Id:                     cluster.Id,
		ProvisioningId:         *cluster.ProvisioningId,
		ProvisioningCredential: decrypted,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (ca *ClusterApplication) OpenTTY(
	ctx context.Context,
	userId int64,
	clusterId int64,
	stdin io.Reader,
	stdout io.Writer,
	sizes <-chan port.TerminalSize,
) error {
	cluster, err := ca.clusterRepository.GetUserCluster(ctx, userId, clusterId)
	if err != nil {
		return err
	}

	if cluster.PrivateKey == nil {
		return fmt.Errorf("cluster private key is not set")
	}
	if cluster.IPv4Address == nil {
		return fmt.Errorf("cluster ipv4 address is not set")
	}

	// The private key was first base64 encoded and then encrypted
	decryptedPrivateKey, err := ca.crypto.Decrypt(*cluster.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %v", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(decryptedPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %v", err)
	}

	pemBytes, err := ca.crypto.EncodePrivateKeyToPEM(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to encode private key to PEM: %v", err)
	}

	return ca.grpcProvisionerClient.OpenTTY(ctx, *cluster.IPv4Address, cluster.User, pemBytes, stdin, stdout, sizes)
}

func (ca *ClusterApplication) StreamProvisioningLogs(
	ctx context.Context,
	userId int64,
	clusterId int64,
	w io.Writer,
) error {
	pr, pw := io.Pipe()
	streamCtx, cancelStream := context.WithCancel(ctx)
	defer cancelStream()

	errCh := make(chan error, 1)
	go func() {
		errCh <- ca.grpcProvisionerClient.StreamProvisioningLogs(streamCtx, clusterId, pw)
		_ = pw.Close()
	}()

	logs, err := ca.clusterRepository.GetUserClusterProvisioningLogs(ctx, userId, clusterId)
	if err != nil {
		cancelStream()
		_ = pr.Close()
		<-errCh
		return err
	}

	if logs != nil && *logs != "" {
		cancelStream()
		_ = pr.Close()
		<-errCh
		_, werr := io.WriteString(w, *logs)
		return werr
	}

	_, copyErr := io.Copy(w, pr)
	grpcErr := <-errCh
	if copyErr != nil {
		return copyErr
	}
	if grpcErr != nil && !errors.Is(grpcErr, context.Canceled) {
		return grpcErr
	}
	return nil
}

func (ca *ClusterApplication) HandleClusterProvisionedSuccess(c *coreValue.ClusterProvisionedSuccess) {
	ctx := context.Background()
	err := ca.clusterRepository.UpdateClusterPulumiStackId(ctx, c.ClusterId, &c.ProvisioningId)
	if err != nil {
		fmt.Printf("failed to persist provisioning id: %v\n", err)
	}

	err = ca.clusterRepository.UpdateClusterIPv4Address(ctx, c.ClusterId, &c.IPv4Address)
	if err != nil {
		fmt.Printf("Failed to persist cluster ip address: %v\n", err)
	}

	encryptedPrivKeyStr, err := ca.crypto.Encrypt(c.PrivateKey)
	if err != nil {
		log.Printf("failed to encrypt private key: %v\n", err)
	}
	err = ca.clusterRepository.UpdateClusterPublicPrivateKey(ctx, c.ClusterId, &c.PublicKey, &encryptedPrivKeyStr)
	if err != nil {
		fmt.Printf("failed to persist cluster public private key: %v\n", err)
	}

	encryptedKubeconfig, err := ca.crypto.Encrypt(c.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to encrypt kubeconfig: %v\n", err)
	}
	err = ca.clusterRepository.UpdateClusterKubeconfig(ctx, c.ClusterId, &encryptedKubeconfig)
	if err != nil {
		fmt.Printf("Failed to persist kubeconfig: %v\n", err)
	}

	err = ca.clusterRepository.UpdateClusterStatus(ctx, c.ClusterId, entity.ClusterRunning)
	if err != nil {
		fmt.Printf("Failed to update cluster status: %v\n", err)
	}

	err = ca.clusterRepository.UpdateClusterLogs(ctx, c.ClusterId, c.Logs)
	if err != nil {
		log.Printf("failed to persist cluster provisioning logs: %v\n", err)
	}

	clusterName := ca.resolveClusterName(ctx, c.ClusterId)
	ca.broadcastClusterNotification(ctx, c.ClusterId, "success", fmt.Sprintf("Cluster %s provisioned successfully", clusterName))
}

func (ca *ClusterApplication) HandleClusterProvisionedFailure(c *coreValue.ClusterProvisionedFailure) {
	ctx := context.Background()
	clusterName := ca.resolveClusterName(ctx, c.ClusterId)
	message := fmt.Sprintf("Failed to provision cluster %s", clusterName)
	if c.Reason != "" {
		message = fmt.Sprintf("%s: %s", message, c.Reason)
	}
	ca.broadcastClusterNotification(ctx, c.ClusterId, "failed", message)
}

func (ca *ClusterApplication) tearDownClusterProjects(ctx context.Context, clusterId int64) error {
	projectIds, err := ca.projectRepository.GetProjectIdsByClusterId(ctx, clusterId)
	if err != nil {
		return err
	}

	for _, projectId := range projectIds {
		envs, err := ca.projectRepository.GetProjectEnvironmentsByProjectId(ctx, projectId)
		if err != nil {
			return err
		}
		for _, env := range envs {
			if err := ca.environmentService.TearDownEnvironmentDeployments(ctx, env); err != nil {
				log.Printf(
					"cluster %d: k8s teardown for environment %d failed (continuing): %v\n",
					clusterId,
					env.Id,
					err,
				)
			}
			if err := ca.deploymentRepository.SoftDeleteDeploymentsByEnvironmentId(ctx, env.Id); err != nil {
				return err
			}
			if err := ca.environmentRepository.DeleteEnvironment(ctx, env.Id); err != nil {
				return err
			}
		}
	}

	return ca.projectRepository.DeleteProjectsByClusterId(ctx, clusterId)
}

func (ca *ClusterApplication) HandleClusterDeletedSuccess(c *coreValue.ClusterDeletedSuccess) {
	ctx := context.Background()

	clusterName := ca.resolveClusterName(ctx, c.ClusterId)

	if err := ca.tearDownClusterProjects(ctx, c.ClusterId); err != nil {
		log.Printf("failed to tear down cluster %d projects: %v\n", c.ClusterId, err)
		ca.broadcastClusterNotification(ctx, c.ClusterId, "failed", fmt.Sprintf("Failed to clean up cluster %s after deletion", clusterName))
		return
	}

	if err := ca.clusterRepository.DeleteCluster(ctx, c.ClusterId); err != nil {
		log.Printf("failed to delete cluster %d from database: %v\n", c.ClusterId, err)
		return
	}

	ca.broadcastClusterNotification(ctx, c.ClusterId, "success", fmt.Sprintf("Cluster %s deleted successfully", clusterName))
}

func (ca *ClusterApplication) HandleClusterDeletedFailure(c *coreValue.ClusterDeletedFailure) {
	ctx := context.Background()
	clusterName := ca.resolveClusterName(ctx, c.ClusterId)
	message := fmt.Sprintf("Failed to delete cluster %s", clusterName)
	if c.Reason != "" {
		message = fmt.Sprintf("%s: %s", message, c.Reason)
	}
	ca.broadcastClusterNotification(ctx, c.ClusterId, "failed", message)
}

func (ca *ClusterApplication) HandleReconcileClusterRequest(req *coreValue.ReconcileClusterRequest) {
	ctx := context.Background()

	cluster, err := ca.clusterRepository.GetCluster(ctx, req.ClusterId)
	if err != nil {
		log.Printf("reconcile skipped: cluster %d: %v\n", req.ClusterId, err)
		return
	}

	if cluster.Status != entity.ClusterRunning {
		log.Printf("reconcile skipped: cluster %d status is %s\n", req.ClusterId, cluster.Status)
		return
	}

	if cluster.ProvisioningId == nil || *cluster.ProvisioningId != req.ProvisioningId {
		log.Printf("reconcile skipped: cluster %d provisioning id mismatch\n", req.ClusterId)
		return
	}

	credential, err := ca.organizationRepository.GetOrganizationProvisioningCredential(
		ctx,
		req.OrganizationId,
		value.HetznerCredential,
	)
	if err != nil {
		log.Printf("reconcile skipped: cluster %d credential: %v\n", req.ClusterId, err)
		return
	}

	decrypted, err := ca.crypto.Decrypt(credential.Secret)
	if err != nil {
		log.Printf("reconcile skipped: cluster %d decrypt credential: %v\n", req.ClusterId, err)
		return
	}

	err = ca.queue.PublishReconcileCluster(&coreValue.ReconcileCluster{
		Id:                     req.ClusterId,
		ProvisioningId:         req.ProvisioningId,
		ProvisioningCredential: decrypted,
	})
	if err != nil {
		log.Printf("failed to publish reconcile cluster for %d: %v\n", req.ClusterId, err)
		return
	}
	log.Printf("reconcile queued for cluster %d (provisioning id %q)\n", req.ClusterId, req.ProvisioningId)
}

func (ca *ClusterApplication) resolveClusterName(ctx context.Context, clusterId int64) string {
	cluster, err := ca.clusterRepository.GetCluster(ctx, clusterId)
	if err != nil || cluster == nil {
		log.Printf("failed to resolve cluster name for %d: %v", clusterId, err)
		return strconv.FormatInt(clusterId, 10)
	}
	return cluster.Name
}

func (ca *ClusterApplication) broadcastClusterNotification(ctx context.Context, clusterId int64, status string, message string) {
	ownerUserId, err := ca.GetClusterOrgOwnerId(ctx, clusterId)
	if err != nil {
		log.Printf("failed to get org owner for cluster %d: %v", clusterId, err)
		return
	}

	ca.userNotificationHub.Broadcast(ownerUserId, &coreValue.ClusterNotification{
		ClusterId: clusterId,
		Status:    status,
		Message:   message,
	})
}

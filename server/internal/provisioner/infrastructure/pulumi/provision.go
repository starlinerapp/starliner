package pulumi

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/crypto/ssh"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/provisioner/domain/port"
)

type DeployParams struct {
	ServerName string
	ServerType value.ServerType
	PublicKey  ed25519.PublicKey
}

func DeployFunc(params DeployParams) pulumi.RunFunc {
	pub, err := ssh.NewPublicKey(params.PublicKey)
	if err != nil {
		return nil
	}

	return func(ctx *pulumi.Context) error {
		sshKeyName := fmt.Sprintf("%s-ssh-key", params.ServerName)
		sshKey, err := hcloud.NewSshKey(ctx, sshKeyName, &hcloud.SshKeyArgs{
			Name:      pulumi.String(sshKeyName),
			PublicKey: pulumi.String(ssh.MarshalAuthorizedKey(pub)),
		})
		if err != nil {
			return err
		}

		s, err := hcloud.NewServer(ctx, params.ServerName, &hcloud.ServerArgs{
			Name:       pulumi.String(params.ServerName),
			Image:      pulumi.String("ubuntu-22.04"),
			ServerType: pulumi.String(params.ServerType),
			Location:   pulumi.String("nbg1"),
			PublicNets: hcloud.ServerPublicNetArray{
				&hcloud.ServerPublicNetArgs{
					Ipv4Enabled: pulumi.Bool(true),
					Ipv6Enabled: pulumi.Bool(true),
				},
			},
			SshKeys: pulumi.StringArray{
				sshKey.ID().ToStringOutput(),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("serverIp", s.Ipv4Address)
		return nil
	}
}

type Provision struct {
	logPublisher port.LogPublisher
}

func NewProvision(
	logPublisher port.LogPublisher,
) port.Provision {
	return &Provision{
		logPublisher: logPublisher,
	}
}

func (p *Provision) ProvisionServer(ctx context.Context, clusterId int64, provisioningCredential string, name string, serverType value.ServerType, publicKey []byte) (provisioningId string, ip string, logs string, err error) {
	stackName := auto.FullyQualifiedStackName("organization", name, uuid.New().String())

	var (
		logBuf strings.Builder
		mu     sync.Mutex
	)
	appendLog := p.logAppender(ctx, clusterId, &logBuf, &mu)

	defer func() {
		logs = logBuf.String()
	}()

	s, err := auto.UpsertStackInlineSource(ctx, stackName, name, DeployFunc(DeployParams{
		ServerName: name,
		ServerType: serverType,
		PublicKey:  publicKey,
	}))
	if err != nil {
		return stackName, "", "", err
	}

	if err := p.configureStack(ctx, s, provisioningCredential); err != nil {
		return stackName, "", "", err
	}

	stream := &inlineLogWriter{appendLog: appendLog}

	_, err = s.Refresh(ctx, optrefresh.ProgressStreams(stream))
	if err != nil {
		return stackName, "", "", err
	}

	res, err := s.Up(ctx, optup.ProgressStreams(stream))
	if err != nil {
		return stackName, "", "", err
	}

	ip, ok := res.Outputs["serverIp"].Value.(string)
	if !ok {
		return stackName, "", "", fmt.Errorf("failed to unmarshall output")
	}

	return stackName, ip, "", nil
}

func (p *Provision) ReconcileServer(ctx context.Context, clusterId int64, provisioningCredential string, provisioningId string) (serverMissing bool, err error) {
	s, stream, err := p.prepareStack(ctx, clusterId, provisioningCredential, provisioningId)
	if err != nil {
		return false, err
	}

	refreshRes, err := s.Refresh(ctx, optrefresh.ProgressStreams(stream))
	if err != nil {
		return false, err
	}

	if refreshRes.Summary.ResourceChanges != nil {
		if deletes := (*refreshRes.Summary.ResourceChanges)["delete"]; deletes > 0 {
			return true, nil
		}
	}

	previewRes, err := s.Preview(ctx, optpreview.ProgressStreams(stream))
	if err != nil {
		return false, err
	}

	return previewRes.ChangeSummary[apitype.OpCreate] > 0, nil
}

func (p *Provision) DestroyServer(ctx context.Context, clusterId int64, provisioningCredential string, provisioningId string) error {
	s, stream, err := p.prepareStack(ctx, clusterId, provisioningCredential, provisioningId)
	if err != nil {
		return err
	}

	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(stream))
	return err
}

func (p *Provision) DeleteServer(ctx context.Context, clusterId int64, provisioningCredential string, provisioningId string) error {
	s, stream, err := p.prepareStack(ctx, clusterId, provisioningCredential, provisioningId)
	if err != nil {
		return err
	}

	_, err = s.Refresh(ctx, optrefresh.ProgressStreams(stream))
	if err != nil {
		return err
	}

	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(stream))
	return err
}

func (p *Provision) prepareStack(
	ctx context.Context,
	clusterId int64,
	provisioningCredential string,
	provisioningId string,
) (auto.Stack, *inlineLogWriter, error) {
	parts := strings.Split(provisioningId, "/")
	if len(parts) < 2 {
		return auto.Stack{}, nil, fmt.Errorf("invalid provisioning id %q", provisioningId)
	}
	projectName := parts[1]

	var (
		logBuf strings.Builder
		mu     sync.Mutex
	)
	appendLog := p.logAppender(ctx, clusterId, &logBuf, &mu)
	stream := &inlineLogWriter{appendLog: appendLog}

	s, err := auto.SelectStackInlineSource(ctx, provisioningId, projectName, DeployFunc(DeployParams{
		ServerName: projectName,
		PublicKey:  nil,
	}))
	if err != nil {
		return auto.Stack{}, nil, err
	}

	if err := p.configureStack(ctx, s, provisioningCredential); err != nil {
		return auto.Stack{}, nil, err
	}

	return s, stream, nil
}

func (p *Provision) configureStack(ctx context.Context, s auto.Stack, provisioningCredential string) error {
	if err := s.SetConfig(ctx, "hcloud:token", auto.ConfigValue{
		Value:  provisioningCredential,
		Secret: true,
	}); err != nil {
		return err
	}

	w := s.Workspace()
	return w.InstallPlugin(ctx, "hcloud", "1.29")
}

func (p *Provision) logAppender(ctx context.Context, clusterId int64, logBuf *strings.Builder, mu *sync.Mutex) func(string) {
	return func(line string) {
		mu.Lock()
		logBuf.WriteString(line)
		mu.Unlock()

		if p.logPublisher != nil {
			if err := p.logPublisher.PublishLogChunk(ctx, clusterId, []byte(line)); err != nil {
				log.Printf("failed to publish log chunk: %v", err)
			}
		}
	}
}

type inlineLogWriter struct {
	appendLog func(line string)
}

func (w *inlineLogWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	w.appendLog(string(p))
	return len(p), nil
}

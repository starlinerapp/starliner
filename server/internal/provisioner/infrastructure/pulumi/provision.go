package pulumi

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"github.com/google/uuid"
	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/crypto/ssh"
	"os"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/provisioner/domain/port"
	"strings"
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
}

func NewProvision() port.Provision {
	return &Provision{}
}

func (p *Provision) ProvisionServer(ctx context.Context, provisioningCredential string, name string, serverType value.ServerType, publicKey []byte) (provisioningId string, ip string, err error) {
	stackName := auto.FullyQualifiedStackName("organization", name, uuid.New().String())

	s, err := auto.UpsertStackInlineSource(ctx, stackName, name, DeployFunc(DeployParams{
		ServerName: name,
		ServerType: serverType,
		PublicKey:  publicKey,
	}))
	if err != nil {
		return stackName, "", err
	}

	err = s.SetConfig(ctx, "hcloud:token", auto.ConfigValue{
		Value:  provisioningCredential,
		Secret: true,
	})
	if err != nil {
		return stackName, "", err
	}

	w := s.Workspace()
	err = w.InstallPlugin(ctx, "hcloud", "1.29")
	if err != nil {
		return stackName, "", err
	}

	_, err = s.Refresh(ctx)
	if err != nil {
		return stackName, "", err
	}

	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		return stackName, "", err
	}

	ip, ok := res.Outputs["serverIp"].Value.(string)
	if !ok {
		return stackName, "", fmt.Errorf("failed to unmarshall output")
	}

	return stackName, ip, nil
}

func (p *Provision) DeleteServer(ctx context.Context, provisioningCredential string, provisioningId string) error {
	parts := strings.Split(provisioningId, "/")
	projectName := parts[1]

	s, err := auto.SelectStackInlineSource(ctx, provisioningId, projectName, DeployFunc(DeployParams{
		ServerName: projectName,
		PublicKey:  nil,
	}))
	if err != nil {
		return err
	}

	err = s.SetConfig(ctx, "hcloud:token", auto.ConfigValue{
		Value:  provisioningCredential,
		Secret: true,
	})
	if err != nil {
		return err
	}

	w := s.Workspace()
	err = w.InstallPlugin(ctx, "hcloud", "1.29")
	if err != nil {
		return err
	}

	_, err = s.Refresh(ctx)
	if err != nil {
		return err
	}

	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)
	_, err = s.Destroy(ctx, stdoutStreamer)
	if err != nil {
		return err
	}

	return nil
}

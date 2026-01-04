package provisioner

import (
	"crypto/ed25519"
	"fmt"
	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/crypto/ssh"
)

type DeployParams struct {
	ServerName string
	publicKey  ed25519.PublicKey
}

func deployFunc(params DeployParams) pulumi.RunFunc {
	pub, err := ssh.NewPublicKey(params.publicKey)
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
			ServerType: pulumi.String("cx23"),
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

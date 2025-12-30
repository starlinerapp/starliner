package cluster

import (
	"github.com/google/uuid"
	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployFunc(ctx *pulumi.Context) error {
	serverName := uuid.New().String()

	s, err := hcloud.NewServer(ctx, serverName, &hcloud.ServerArgs{
		Name:       pulumi.String(serverName),
		Image:      pulumi.String("ubuntu-22.04"),
		ServerType: pulumi.String("cx23"),
		Location:   pulumi.String("nbg1"),
		PublicNets: hcloud.ServerPublicNetArray{
			&hcloud.ServerPublicNetArgs{
				Ipv4Enabled: pulumi.Bool(true),
				Ipv6Enabled: pulumi.Bool(true),
			},
		},
	})
	if err != nil {
		return err
	}

	ctx.Export("serverIp", s.Ipv4Address)
	return nil
}

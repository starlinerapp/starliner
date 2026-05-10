package application

import (
	"fmt"

	"starliner.app/cli/internal/domain/port"
)

type DemoApplication struct {
	k3dClient port.K3dClient
	wgClient  port.WgClient
}

func NewDemoApplication(
	k3dClient port.K3dClient,
	wgClient port.WgClient,
) *DemoApplication {
	return &DemoApplication{
		k3dClient: k3dClient,
		wgClient:  wgClient,
	}
}

func (da *DemoApplication) Run() error {
	if err := da.k3dClient.Install(); err != nil {
		return err
	}
	if err := da.wgClient.Install(); err != nil {
		return err
	}
	keyPair, err := da.wgClient.GenerateKeyPair()
	if err != nil {
		return err
	}

	fmt.Printf("public key: %s\n", keyPair.PublicKey)
	fmt.Printf("private key: %s\n", keyPair.PrivateKey)

	if err := da.k3dClient.Start(); err != nil {
		return err
	}

	return nil
}

func (da *DemoApplication) Stop() error {
	if err := da.k3dClient.Stop(); err != nil {
		return err
	}

	return nil
}

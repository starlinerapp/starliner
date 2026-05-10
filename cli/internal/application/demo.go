package application

import "starliner.app/cli/internal/domain/port"

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

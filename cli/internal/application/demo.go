package application

import (
	"fmt"
	"os"
	"os/exec"
)

type DemoApplication struct {
}

func NewDemoApplication() *DemoApplication {
	return &DemoApplication{}
}

func (da *DemoApplication) Run() error {
	if err := da.installK3d(); err != nil {
		return err
	}

	if err := da.startK3d(); err != nil {
		return err
	}

	return nil
}

func (da *DemoApplication) installK3d() error {
	fmt.Println("Installing k3d...")
	cmd := exec.Command(
		"bash",
		"-c",
		"wget -q -O - https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install k3d: %w", err)
	}

	fmt.Println("k3d installed successfully")

	return nil
}

func (da *DemoApplication) startK3d() error {
	fmt.Println("Starting k3d cluster...")

	cmd := exec.Command(
		"k3d",
		"cluster",
		"create",
		"local",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start k3d cluster: %w", err)
	}

	fmt.Println("k3d cluster started successfully")

	return nil
}

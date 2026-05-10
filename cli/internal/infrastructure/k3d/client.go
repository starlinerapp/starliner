package k3d

import (
	"fmt"
	"os"
	"os/exec"

	"starliner.app/cli/internal/domain/port"
)

type Client struct {
}

func NewClient() port.K3dClient {
	return &Client{}
}

func (c *Client) Install() error {
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

func (c *Client) Start() error {
	fmt.Println("Starting k3d cluster...")

	cmd := exec.Command(
		"k3d",
		"cluster",
		"create",
		"starliner-demo",
		"--api-port=16444",
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

func (c *Client) Stop() error {
	fmt.Println("Stopping k3d cluster...")
	cmd := exec.Command(
		"k3d",
		"cluster",
		"delete",
		"starliner-demo",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete k3d cluster: %w", err)
	}

	fmt.Println("k3d cluster deleted successfully")

	return nil
}

package wg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"starliner.app/cli/internal/domain/port"
)

type Client struct {
}

func NewClient() port.WgClient {
	return &Client{}
}

func (c *Client) Install() error {
	fmt.Println("Installing WireGuard...")

	switch runtime.GOOS {
	case "darwin":
		return installOnMacOS()
	case "linux":
		return installOnLinux()
	case "windows":
		return nil
	default:
		panic("unsupported platform")
	}
}

func (c *Client) GenerateKeyPair() (keyPair *port.KeyPair, err error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	return &port.KeyPair{
		PublicKey:  privateKey.PublicKey().String(),
		PrivateKey: privateKey.String(),
	}, nil
}

func installOnMacOS() error {
	cmd := exec.Command(
		"brew",
		"install",
		"wireguard-tools",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install wireguard: %w", err)
	}

	return nil
}

func installOnLinux() error {
	var cmd *exec.Cmd
	switch {
	case commandExists("apt-get"):
		cmd = exec.Command(
			"sudo",
			"apt-get",
			"install",
			"-y",
			"wireguard",
		)
	case commandExists("dnf"):
		cmd = exec.Command(
			"sudo",
			"dnf",
			"install",
			"-y",
			"wireguard-tools",
		)

	case commandExists("yum"):
		cmd = exec.Command(
			"sudo",
			"yum",
			"install",
			"-y",
			"wireguard-tools",
		)

	case commandExists("apk"):
		cmd = exec.Command(
			"sudo",
			"apk",
			"add",
			"wireguard-tools",
		)
	default:
		return fmt.Errorf("unsupported linux package manager")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install wireguard: %w", err)
	}

	return nil
}

func (c *Client) installWindows() error {
	var cmd *exec.Cmd

	switch {
	case commandExists("winget"):
		cmd = exec.Command(
			"winget",
			"install",
			"--id",
			"WireGuard.WireGuard",
			"-e",
			"--silent",
		)

	case commandExists("choco"):
		cmd = exec.Command(
			"choco",
			"install",
			"wireguard",
			"-y",
		)

	default:
		return fmt.Errorf("no supported Windows package manager found")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install wireguard: %w", err)
	}

	return nil
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

package ansible

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/provisioner/domain/port"
)

type K3sPlaybookOutput struct {
	Plays []struct {
		Tasks []struct {
			Hosts map[string]struct {
				Content string `json:"content,omitempty"`
			} `json:"hosts"`
		} `json:"tasks"`
	} `json:"plays"`
}

type Install struct {
	crypto corePort.Crypto
}

func NewInstall(crypto corePort.Crypto) port.Install {
	return &Install{
		crypto: crypto,
	}
}

func (i *Install) InstallK3s(ip string, privateKey []byte) (kubeconfig string, err error) {
	tmpPlaybook, err := os.CreateTemp("", "ansible-*.yml")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPlaybook.Name())

	_, err = tmpPlaybook.Write([]byte(K3sPlaybook))
	if err != nil {
		return "", err
	}
	err = tmpPlaybook.Close()
	if err != nil {
		return "", err
	}

	tmpPrivateKey, err := os.CreateTemp("", "private-key-*.pem")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPrivateKey.Name())

	pemBytes, err := i.crypto.EncodePrivateKeyToPEM(privateKey)
	if err != nil {
		return "", err
	}
	_, err = tmpPrivateKey.Write(pemBytes)
	if err != nil {
		return "", err
	}
	err = tmpPrivateKey.Close()
	if err != nil {
		return "", err
	}

	args := []string{
		tmpPlaybook.Name(),
		"-i", fmt.Sprintf("%s,", ip),
		"-u", "root",
		"--private-key", tmpPrivateKey.Name(),
		"-e", "ansible_ssh_common_args='-o StrictHostKeyChecking=no'",
		"-e", fmt.Sprintf("target_host=%s", ip),
	}

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Env = append(
		cmd.Env,
		"ANSIBLE_STDOUT_CALLBACK=json",
		"ANSIBLE_DEPRECATION_WARNINGS=False",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	var ansibleData K3sPlaybookOutput
	if err := json.Unmarshal(out, &ansibleData); err != nil {
		return "", err
	}

	var kubeconfigBase64 string
	for _, play := range ansibleData.Plays {
		for _, task := range play.Tasks {
			for _, host := range task.Hosts {
				if host.Content != "" {
					kubeconfigBase64 = host.Content
					break
				}
			}
		}
	}
	kubeconfigDecoded, err := base64.StdEncoding.DecodeString(kubeconfigBase64)
	if err != nil {
		fmt.Printf("failed to decode kubeconfig: %v\n", err)
	}

	kubeconfig = strings.ReplaceAll(
		string(kubeconfigDecoded),
		"https://127.0.0.1:",
		fmt.Sprintf("https://%s:", ip),
	)

	return kubeconfig, nil
}
